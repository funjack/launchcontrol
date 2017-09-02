# ##### BEGIN BSD LICENSE BLOCK #####
#
# Copyright 2017 Funjack
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
#
# 1. Redistributions of source code must retain the above copyright notice, this
# list of conditions and the following disclaimer.
#
# 2. Redistributions in binary form must reproduce the above copyright notice,
# this list of conditions and the following disclaimer in the documentation
# and/or other materials provided with the distribution.
#
# 3. Neither the name of the copyright holder nor the names of its contributors
# may be used to endorse or promote products derived from this software without
# specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
# ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
# WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
# SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
# CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
# OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
#
# ##### END BSD LICENSE BLOCK #####

bl_info = {
    "name": "Funscripting Addon",
    "author": "Funjack",
    "version": (0, 0, 6),
    "location": "Sequencer",
    "description": "Script Launch haptics data and export as Funscript.",
    "category": "Sequencer",
}

import bpy
import json

addon_keymaps = []

class FunscriptPanel(bpy.types.Panel):
    """Funscript UI panel.

    Funscript UI panel added to the sequencer.
    """
    bl_label = "Funscript"
    bl_space_type = "SEQUENCE_EDITOR"
    bl_region_type = "UI"

    @classmethod
    def poll(cls, context):
        return context.selected_sequences is not None \
                and len(context.selected_sequences) == 1

    def limitinfo(self, context):
        """Labels with hints of the limitations"""
        scene = context.scene
        seq = context.selected_sequences[0]
        keyframes = launch_keyframes(seq.name)
        layout = self.layout
        row = layout.row(align=True)
        col = row.column(align=True)
        last = {"frame":0, "value":0}
        if keyframes is not None:
            for kf in reversed(keyframes):
                frame = kf.co[0]
                value = kf.co[1]
                if frame > scene.frame_current:
                    continue
                if frame < scene.frame_current:
                    last = {"frame":frame, "value":value}
                    break
        interval = frame_to_ms(scene.frame_current) - frame_to_ms(last["frame"])
        icon = "FILE_TICK" if interval > 100 or last["frame"] == 0 else "ERROR"
        if interval > 1000:
            icon = "TIME"
        mindist = launch_distance(20, interval)
        maxdist = launch_distance(80, interval)
        col.label(text="Previous: %d" % last["value"])
        col = row.column(align=True)
        col.label("Slowest: %d" % mindist)
        row = layout.row(align=True)
        col = row.column(align=True)
        col.label(text="Interval: %d ms" % interval, icon=icon)
        col = row.column(align=True)
        col.label("Fastest: %d" % maxdist)

    def draw(self, context):
        self.limitinfo(context)
        layout = self.layout
        for x in [0, 10, 40, 70, 100]:
            row = layout.row(align=True)
            row.alignment = 'EXPAND'
            if x == 0 or x == 100:
                row.operator("funscript.position", text=str(x)).launchPosition=x
            else:
                for i in range(x,x+30,10):
                    row.operator("funscript.position", text=str(i)).launchPosition=i
        layout.label(text="Generate")
        row = layout.row(align=True)
        row.alignment = 'EXPAND'
        row.operator("funscript.repeat")
        row.operator("funscript.fill")
        layout.label(text="Import funscript")
        row = layout.row(align=True)
        row.alignment = 'EXPAND'
        row.operator("funscript.import")
        layout.label(text="Export funscript")
        row = layout.row(align=True)
        row.alignment = 'EXPAND'
        row.operator("funscript.export")

class FunscriptPositionButton(bpy.types.Operator):
    """Position input button.

    Button that inserts a Launch position in the currently selected Sequence.
    """
    bl_idname = "funscript.position"
    bl_label = "Position"
    bl_options = {'REGISTER', 'UNDO'}
    launchPosition = bpy.props.IntProperty()

    def execute(self, context):
        print("inserting: %d" % self.launchPosition)
        scene = context.scene
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        insert_position(seq, self.launchPosition, scene.frame_current)
        scene.frame_set(scene.frame_current)
        return{'FINISHED'}

class FunscriptRepeatButton(bpy.types.Operator):
    """Repeat last stroke button.

    Button that will repeat the last stroke on the selected sequence.
    """
    bl_idname = "funscript.repeat"
    bl_label = "Repeat stroke"
    bl_options = {'REGISTER', 'UNDO'}

    def execute(self, context):
        scene = context.scene
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        lastframe = repeat_stroke(seq, scene.frame_current)
        scene.frame_set(lastframe)
        return{'FINISHED'}

class FunscriptFillButton(bpy.types.Operator):
    """Fill last stroke button.

    Button that will repeat the last stroke on the selected sequence until the
    current frame is reached.
    """
    bl_idname = "funscript.fill"
    bl_label = "Fill stroke"
    bl_options = {'REGISTER', 'UNDO'}

    def execute(self, context):
        scene = context.scene
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        lastframe = repeat_fill_stroke(seq, scene.frame_current)
        scene.frame_set(lastframe)
        return{'FINISHED'}

class FunscriptExport(bpy.types.Operator):
    """Export as Funscript file button.

    Button that exports all Launch position keyframes in the sequences as
    Funscript file.
    """
    bl_idname = "funscript.export"
    bl_label = "Export as Funscript"
    filepath = bpy.props.StringProperty(subtype='FILE_PATH')
    inverted = bpy.props.BoolProperty(name="inverted",
        description="Flip up and down positions", default=False)

    def execute(self, context):
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        keyframes = launch_keyframes(seq.name)
        script = create_funscript(keyframes, self.inverted)
        with open(self.filepath, 'w') as outfile:
            json.dump(script, outfile)
        return {'FINISHED'}

    def draw(self,context):
        layout = self.layout
        layout.prop(self, "inverted", text="Inverted")

    def invoke(self, context, event):
        context.window_manager.fileselect_add(self)
        return {'RUNNING_MODAL'}

class FunscriptImport(bpy.types.Operator):
    """Import as Funscript file button.

    Button that imports Launch position keyframes in the sequences from a
    Funscript file.
    """
    bl_idname = "funscript.import"
    bl_label = "Import Funscript on frame"
    bl_options = {'REGISTER', 'UNDO'}
    filepath = bpy.props.StringProperty(subtype='FILE_PATH')

    def execute(self, context):
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        with open(self.filepath) as infile:
            fs = json.load(infile)
            if not "actions" in fs:
                self.report({'ERROR_INVALID_INPUT'}, "Input is not valid funscript.")
                return{'CANCELLED'}
            insert_actions(seq, fs["actions"], context.scene.frame_current)
        return {'FINISHED'}

    def invoke(self, context, event):
        context.window_manager.fileselect_add(self)
        return {'RUNNING_MODAL'}

def insert_actions(seq, actions, offset=0):
    """Insert into seq positions from the actions dict"""
    for a in actions:
        frame = ms_to_frame(a["at"]) + offset
        insert_position(seq, a["pos"], frame)

def insert_stroke(seq, stroke, offset=0):
    """Insert into seq positions from the stroke dict"""
    frame = offset
    for i, p in enumerate(stroke):
        frame = p["frame"] + offset
        if i == 0:
            # Do not override a keyframe at the start of the insert position.
            if launch_keyframe(seq.name, frame) is not None:
                continue
        # XXX Maybe we should also remove any keyframes that already exists in
        # the frame window of the new stroke?
        insert_position(seq, p["value"], frame)
    return frame

def insert_position(seq, position, frame):
    """Inserts in seq a keyframe with value position on frame."""
    seq["launch"] = position
    seq.keyframe_insert(data_path='["launch"]', frame=frame)

def repeat_stroke(seq, frame_current):
    """Repeat the last stroke on the current frame"""
    stroke = last_stroke(seq, frame_current)
    if stroke is None or len(stroke) < 3:
        return
    return insert_stroke(seq, stroke, frame_current)

def repeat_fill_stroke(seq, frame_end):
    """Fill the the last stroke before end_frame until end_frame."""
    stroke = last_stroke(seq, frame_end)
    if stroke is None or len(stroke) < 3:
        return
    frame = frame_end
    keyframes = launch_keyframes(seq.name)
    for kf in reversed(keyframes):
        frame = kf.co[0]
        if frame > frame_end:
            continue
        if frame <= frame_end:
            break
    return fill_stroke(seq, stroke, frame, frame_end)

def fill_stroke(seq, stroke, frame_start, frame_end):
    """Fill between frame_start and frame_end with stroke."""
    if stroke is None or len(stroke) < 3:
        return
    frame = frame_start
    while frame + stroke[-1]["frame"] < frame_end:
        frame = insert_stroke(seq, stroke, frame)
    return frame

def last_stroke(seq, since_frame):
    """Returns the last stroke since frame."""
    keyframes = launch_keyframes(seq.name)
    stroke = []
    for kf in reversed(keyframes):
        frame = kf.co[0]
        value = kf.co[1]
        if frame > since_frame:
            continue
        if frame <= since_frame:
            stroke.append({"frame":frame, "value":value})
        if len(stroke) == 3:
            break
    if len(stroke) < 3:
        return
    startframe = stroke[1]["frame"] - stroke[2]["frame"]
    endframe = stroke[0]["frame"] - stroke[2]["frame"]
    return [ {"frame": 0, "value": stroke[2]["value"] },
             {"frame": startframe, "value": stroke[1]["value"] },
             {"frame": endframe, "value": stroke[0]["value"] } ]

def create_funscript(keyframes, inverted):
    """Create Funscript from keyframes."""
    script = []
    for kf in keyframes:
        time = frame_to_ms(int(kf.co[0]))
        if time < 0:
            continue
        value = int(kf.co[1])
        if value < 0:
            value = 0
        elif value > 100:
            value = 100
        script.append({"at": time, "pos": value})
    return {"version": "1.0", "inverted": inverted, "range": 100, "actions": script}

def launch_keyframes(name):
    """Return all keyframes from all actions fcurves in prop 'launch'."""
    for a in bpy.data.actions:
        for f in a.fcurves:
            if f.data_path.endswith('["%s"]["launch"]' % name):
                return f.keyframe_points

def launch_keyframe(name, frame):
    """Returns the keyframe value at frame."""
    keyframes = launch_keyframes(name)
    for kf in keyframes:
        if kf.co[0] == frame:
            return kf.co[1]
    return None

def frame_to_ms(frame):
    """Returns time position in milliseconds for the given frame number."""
    scene = bpy.context.scene
    fps = scene.render.fps
    fps_base = scene.render.fps_base
    return round((frame-1)/fps*fps_base*1000)

def ms_to_frame(time):
    """Returns frame number for the given time position in milliseconds."""
    scene = bpy.context.scene
    fps = scene.render.fps
    fps_base = scene.render.fps_base
    return round(time/1000/fps_base*fps+1)

def launch_distance(speed, duration):
    """Returns the launch movement distance for given speed and time in ms."""  
    if speed <= 0 or duration <= 0:
        return 0
    time = pow(speed/25000, -0.95)
    diff = time - duration
    return 90 - int(diff/time*90)

def register():
    bpy.utils.register_class(FunscriptPositionButton)
    bpy.utils.register_class(FunscriptRepeatButton)
    bpy.utils.register_class(FunscriptFillButton)
    bpy.utils.register_class(FunscriptExport)
    bpy.utils.register_class(FunscriptImport)
    bpy.utils.register_class(FunscriptPanel)

    # handle the keymap
    wm = bpy.context.window_manager
    kc = wm.keyconfigs.addon
    if kc:
        km = wm.keyconfigs.addon.keymaps.new(name='Sequencer', space_type='SEQUENCE_EDITOR')
        kmi = km.keymap_items.new(FunscriptFillButton.bl_idname, 'EQUAL', 'PRESS')
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptRepeatButton.bl_idname, 'ACCENT_GRAVE', 'PRESS')
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'ZERO', 'PRESS')
        kmi.properties.launchPosition = 0
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_0', 'PRESS')
        kmi.properties.launchPosition = 0
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'ONE', 'PRESS')
        kmi.properties.launchPosition = 10
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_1', 'PRESS')
        kmi.properties.launchPosition = 10
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'TWO', 'PRESS')
        kmi.properties.launchPosition = 20
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_2', 'PRESS')
        kmi.properties.launchPosition = 20
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'THREE', 'PRESS')
        kmi.properties.launchPosition = 30
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_3', 'PRESS')
        kmi.properties.launchPosition = 30
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'FOUR', 'PRESS')
        kmi.properties.launchPosition = 40
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_4', 'PRESS')
        kmi.properties.launchPosition = 40
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'FIVE', 'PRESS')
        kmi.properties.launchPosition = 50
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_5', 'PRESS')
        kmi.properties.launchPosition = 50
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'SIX', 'PRESS')
        kmi.properties.launchPosition = 60
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_6', 'PRESS')
        kmi.properties.launchPosition = 60
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'SEVEN', 'PRESS')
        kmi.properties.launchPosition = 70
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_7', 'PRESS')
        kmi.properties.launchPosition = 70
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'EIGHT', 'PRESS')
        kmi.properties.launchPosition = 80
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_8', 'PRESS')
        kmi.properties.launchPosition = 80
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NINE', 'PRESS')
        kmi.properties.launchPosition = 90
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_9', 'PRESS')
        kmi.properties.launchPosition = 90
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'MINUS', 'PRESS')
        kmi.properties.launchPosition = 100
        addon_keymaps.append((km, kmi))
        kmi = km.keymap_items.new(FunscriptPositionButton.bl_idname, 'NUMPAD_PERIOD', 'PRESS')
        kmi.properties.launchPosition = 100
        addon_keymaps.append((km, kmi))

def unregister():
    for km, kmi in addon_keymaps:
        km.keymap_items.remove(kmi)
    addon_keymaps.clear()

    bpy.utils.unregister_class(FunscriptPanel)
    bpy.utils.unregister_class(FunscriptImport)
    bpy.utils.unregister_class(FunscriptExport)
    bpy.utils.unregister_class(FunscriptFillButton)
    bpy.utils.unregister_class(FunscriptRepeatButton)
    bpy.utils.unregister_class(FunscriptPositionButton)

if __name__ == "__main__":
    register()
