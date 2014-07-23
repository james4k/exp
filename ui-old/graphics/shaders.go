package graphics

import (
	"j4k.co/gfx"
)

var vs gfx.VertexShader = `
#version 150

uniform mat4 WorldViewProjectionM;

in vec3 Position;
in vec4 Color;
in vec2 UV;
//in vec3 Normal;

out vec2 uv;
out vec4 color;

void main() {
	/*
	vec3 lightdir = vec3(0.0, 0.7, -0.7);
	//vec4 worldNormal = WorldViewProjectionM * vec4(Normal, 1.0);
	vec4 worldNormal = vec4(lightdir, 1.0);
	float costheta = clamp(dot(worldNormal.xyz, lightdir), 0.0, 1.0);
	vec3 light = vec3(0.4) + vec3(1.0, 1.0, 1.0) * costheta;
	*/
	color = Color;
	//uv = vec2(UV.x, 1.0 - UV.y);
	uv = UV;
	gl_Position = WorldViewProjectionM * vec4(Position, 1.0);
}`

var fs gfx.FragmentShader = `
#version 150
//#extension GL_ARB_explicit_attrib_location : enable

uniform sampler2D Diffuse;

in vec2 uv;
in vec4 color;

out vec4 FinalColor;

void main() {
	vec4 diff = texture(Diffuse, uv);
	// "gamma correct" bright colors - this is not 'correct' at all but certainly
	// increases legibility for white-on-black.
	//float mask = pow(diff.r, 1.0 - length(color.rgb)*0.15);
	//FinalColor = vec4(color.rgb, mask * color.a);
	FinalColor = color;
}`
