import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import ts from 'typescript';

const source=await readFile(new URL('../src/road-label.ts',import.meta.url),'utf8');
const compiled=ts.transpileModule(source,{compilerOptions:{target:ts.ScriptTarget.ES2022,module:ts.ModuleKind.ESNext}}).outputText;
const {RoadLabel}=await import('data:text/javascript;base64,'+Buffer.from(compiled).toString('base64'));
class Point{constructor(x,y){this.x=x;this.y=y}clone(){return new Point(this.x,this.y)}}
const feature={props:{name:'Hook Lane'}};
const scratch={measureText:t=>({width:t.length*6})};
const style=new RoadLabel(11,'#4c5865','#ffffff');
const place=points=>style.place({scratch},[points.map(p=>new Point(...p))],feature);
const curve=[[0,0],[35,0],[70,3],[105,9],[140,18],[180,32]];
for(const points of [curve,[...curve].reverse()]){
 const labels=place(points);assert(labels?.length,'A gently bending road must retain its name');
 let text='',angles=[];
 labels[0].draw({save(){},restore(){},translate(){},rotate:a=>angles.push(a),strokeText(){},fillText:t=>text+=t});
 assert.equal(text,'Hook Lane');
 assert(angles.every(a=>Math.abs(a)<Math.PI/2),'Reversed geometry must still draw upright text');
 assert(new Set(angles).size>1,'Glyphs must follow the bend rather than use a single straight baseline');
 assert(labels[0].bboxes.every(b=>Object.values(b).every(Number.isFinite)&&b.maxX>b.minX&&b.maxY>b.minY));
}
assert.equal(place([[0,0],[20,0]]),undefined,'Do not squeeze a name onto a short road');
assert.equal(place([[0,0],[0,0],[0,0]]),undefined,'Degenerate geometry must not create invalid labels');
assert.equal(place([[0,0],[40,0],[40,30],[0,30]]),undefined,'Do not wrap names around a hairpin');
console.log('Road-label geometry checks passed: curved/reversed roads, short lines, degenerate lines and sharp turns.');
