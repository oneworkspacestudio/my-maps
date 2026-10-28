import type {Label,LabelSymbolizer,Layout,Feature} from 'protomaps-leaflet';

type Geometry=Parameters<LabelSymbolizer['place']>[1];

// Follow the road itself, rather than requiring a long, almost straight segment.
// This keeps names on gently bending streets readable at neighbourhood zooms.
export class RoadLabel implements LabelSymbolizer{
 constructor(private size:number,private color:string,private halo:string){}
 place(layout:Layout,geometry:Geometry,feature:Feature):Label[]|undefined{
  const name=String(feature.props['name:en']||feature.props.name||'');
  if(!name||name.length>48)return;
  const font=`400 ${this.size}px Arial, sans-serif`,ctx=layout.scratch;
  ctx.font=font;
  const chars=Array.from(name),widths=chars.map(c=>ctx.measureText(c).width);
  const width=widths.reduce((sum,w)=>sum+w,0),labels:Label[]=[];
  for(const original of geometry){
   if(original.length<2)continue;
   // Reverse a copy so text stays upright without changing shared tile geometry.
   const points=original[original.length-1].x<original[0].x?[...original].reverse():original;
   const lengths=[0];
   for(let i=1;i<points.length;i++)lengths.push(lengths[i-1]+Math.hypot(points[i].x-points[i-1].x,points[i].y-points[i-1].y));
   const length=lengths[lengths.length-1];
   if(length<width+18)continue;
   const at=(distance:number)=>{
    let i=1;while(i<lengths.length-1&&lengths[i]<distance)i++;
    const a=points[i-1],b=points[i],span=lengths[i]-lengths[i-1];
    const t=span?(distance-lengths[i-1])/span:0;
    return {x:a.x+(b.x-a.x)*t,y:a.y+(b.y-a.y)*t,angle:Math.atan2(b.y-a.y,b.x-a.x)};
   };
   const count=Math.max(1,Math.floor(length/280));
   for(let slot=0;slot<count;slot++){
    const start=(slot+.5)*length/count-width/2;
    let distance=start,previous:number|undefined,turn=0,valid=true;
    const glyphs=chars.map((char,i)=>{
     const p=at(distance+widths[i]/2);distance+=widths[i];
     if(previous!==undefined){const change=Math.abs(Math.atan2(Math.sin(p.angle-previous),Math.cos(p.angle-previous)));turn+=change;if(change>Math.PI/5)valid=false}
     previous=p.angle;return {...p,char,width:widths[i]};
    });
    // Do not bend names around a hairpin, roundabout or sharp junction.
    if(!valid||turn>Math.PI/2)continue;
    const anchor=points[0].clone(),middle=at(start+width/2);anchor.x=middle.x;anchor.y=middle.y;
    labels.push({anchor,deduplicationKey:name,deduplicationDistance:240,
     bboxes:glyphs.filter(g=>g.char!==' ').map(g=>{
      const w=Math.abs(Math.cos(g.angle))*g.width/2+Math.abs(Math.sin(g.angle))*this.size/2+2;
      const h=Math.abs(Math.sin(g.angle))*g.width/2+Math.abs(Math.cos(g.angle))*this.size/2+2;
      return {minX:g.x-w,maxX:g.x+w,minY:g.y-h,maxY:g.y+h};
     }),
     draw:canvas=>{
      canvas.font=font;canvas.fillStyle=this.color;canvas.strokeStyle=this.halo;
      canvas.lineWidth=2.5;canvas.lineJoin='round';canvas.globalAlpha=1;
      for(const g of glyphs){
       canvas.save();canvas.translate(g.x-anchor.x,g.y-anchor.y);canvas.rotate(g.angle);
       canvas.strokeText(g.char,-g.width/2,this.size*.34);canvas.fillText(g.char,-g.width/2,this.size*.34);canvas.restore();
      }
     },
    });
   }
  }
  return labels.length?labels:undefined;
 }
}
