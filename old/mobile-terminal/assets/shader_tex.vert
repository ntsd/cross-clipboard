#version 100
attribute vec3 position;
attribute vec2 textureCoord;
uniform mat4 MVP;
varying vec2 vTextureCoord;
void main() {
	vTextureCoord = textureCoord;
	gl_Position = MVP * vec4(position,1.0);
}
