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
    "version": (0, 0, 3),
    "location": "Sequencer",
    "description": "Script Launch haptics data and export as Funscript.",
    "category": "Sequencer",
}

import bpy
import json

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

    def draw(self, context):
        layout = self.layout

        for x in [0, 10, 40, 70, 100]:
            row = layout.row(align=True)
            row.alignment = 'EXPAND'
            if x == 0 or x == 100:
                row.operator("funscript.position", text=str(x)).launchPosition=x
            else:
                for i in range(x,x+30,10):
                    row.operator("funscript.position", text=str(i)).launchPosition=i

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
    launchPosition = bpy.props.IntProperty()

    def execute(self, context):
        scene = context.scene
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        insert_position(seq, self.launchPosition, scene.frame_current)
        return{'FINISHED'}

class FunscriptExport(bpy.types.Operator):
    """Export as Funscript file button.

    Button that exports all Launch position keyframes in the sequences as
    Funscript file.
    """
    bl_idname = "funscript.export"
    bl_label = "Export as Funscript"
    filepath = bpy.props.StringProperty(subtype='FILE_PATH')

    def execute(self, context):
        if len(context.selected_sequences) < 1:
            self.report({'ERROR_INVALID_CONTEXT'}, "No sequence selected.")
            return{'CANCELLED'}
        seq = context.selected_sequences[0]
        keyframes = launch_keyframes(seq.name)
        script = create_funscript(keyframes)
        with open(self.filepath, 'w') as outfile:
            json.dump(script, outfile)
        return {'FINISHED'}

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
    """Insert into seq positions from actions dict"""
    for a in actions:
        frame = ms_to_frame(a["at"]) + offset
        insert_position(seq, a["pos"], frame)

def insert_position(seq, position, frame):
    """Inserts in seq a keyframe with value position on frame."""
    seq["launch"] = position
    seq.keyframe_insert(data_path='["launch"]', frame=frame)

def create_funscript(keyframes):
    """Create Funscript from keyframes."""
    script = []
    for kf in keyframes:
        time = frame_to_ms(int(kf.co[0]))
        value = int(kf.co[1])
        script.append({"at": time, "pos": value})
    return {"version": "1.0", "inverted": True, "range": 100, "actions": script}

def launch_keyframes(name):
    """Return all keyframes from all actions fcurves in prop 'launch'."""
    for a in bpy.data.actions:
        for f in a.fcurves:
            if f.data_path.endswith('["%s"]["launch"]' % name):
                return f.keyframe_points

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


def register():
    bpy.utils.register_class(FunscriptPositionButton)
    bpy.utils.register_class(FunscriptExport)
    bpy.utils.register_class(FunscriptImport)
    bpy.utils.register_class(FunscriptPanel)

def unregister():
    bpy.utils.unregister_class(FunscriptPositionButton)
    bpy.utils.unregister_class(FunscriptExport)
    bpy.utils.unregister_class(FunscriptImport)
    bpy.utils.unregister_class(FunscriptPanel)

if __name__ == "__main__":
    register()
