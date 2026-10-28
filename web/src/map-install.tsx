import React,{useEffect,useRef,useState} from 'react';

import {MapInstall as Screen} from './generated/MapInstall';
import {DirectionsInstall} from './directions-install';

const size=(n:number)=>n>=1073741824?(n/1073741824).toFixed(2)+' GB':(n/1048576).toFixed(1)+' MB';

const coverageName=(c:any)=>({world:'World',uk:'United Kingdom',europe:'Europe',area:'Local area'}[c.coverage as string]||c.coverage);

const label=(c:any)=>coverageName(c)+' · '+({overview:'Overview',region:'Regional detail',streets:'Street detail'}[c.depth as string]||c.depth)+(c.coverage==='area'?' · '+c.radius+' km radius':'');

export function MapInstall(p:any){

 const [fields,setFields]=useState<any>({coverage:p.config.installed?.choice.coverage||'uk',depth:p.config.installed?.choice.depth||'region',lat:Number(p.centre?.lat??51.505).toFixed(5),lon:Number(p.centre?.lon??-.09).toFixed(5),radius:String(p.config.installed?.choice.radius||10),postcodes:p.config.installed?.choice.postcodes??true,includeDirections:!p.config.directions?.ready,router:''}),[status,setStatus]=useState<any>(null),[plan,setPlan]=useState<any>(null),[directions,setDirections]=useState<any>(p.config.directions),[queueStage,setQueueStage]=useState('');
 const queued=useRef<any>(null),alive=useRef(true);

 useEffect(()=>{alive.current=true;let loading=false;const load=async()=>{if(loading)return;loading=true;try{const [s,d]=await Promise.all([p.request('/api/map-install'),p.request('/api/directions')]);if(!alive.current)return;setStatus(s);setDirections(d);
 const q=queued.current;if(q&&q.stage==='map'&&s.installed?.id!==q.previous&&s.progress?.phase==='ready'&&!s.progress.busy){q.stage='starting';setQueueStage('Checking local directions...');try{const fresh=await p.request('/api/directions/plan',{method:'POST',body:JSON.stringify(q.choice)});if(fresh.downloadBytes>q.bytes)throw Error('The directions download grew. Check its size again below. Your map is ready.');await p.request('/api/directions/start',{method:'POST',body:JSON.stringify({plan:fresh.id})});q.stage='directions';if(alive.current)setQueueStage('Installing local directions...')}catch(e:any){if(alive.current){queued.current=null;setQueueStage('');p.setError(e.message)}}}
 if(q&&q.stage==='map'&&!s.progress.busy&&['cancelled','error'].includes(s.progress.phase)){queued.current=null;setQueueStage('')}
 if(q?.stage==='directions'&&d.installed?.id!==q.previousDir&&d.progress?.phase==='ready'&&!d.progress.busy){queued.current=null;setQueueStage('')}
 if(q?.stage==='directions'&&['error','cancelled'].includes(d.progress?.phase)){queued.current=null;setQueueStage('')}
 }catch(e:any){if(alive.current)p.setError(e.message)}finally{loading=false}};void load();const timer=setInterval(load,1500);return()=>{alive.current=false;clearInterval(timer)}},[]);


 const change=(key:string,value:any)=>{setPlan(null);setFields((f:any)=>{const next={...f,[key]:value};if(key==='coverage'&&value==='world'){next.depth='overview';next.postcodes=false;next.includeDirections=false}if(key==='coverage'&&value==='europe')next.includeDirections=false;if(key==='coverage'&&['uk','area'].includes(value))next.includeDirections=!directions?.ready;if(key==='coverage'&&value!=='area'&&next.depth==='streets')next.depth='region';return next})};

 const action=async(a:string)=>{p.setError('');try{switch(a){

  case'mapBack':case'openMap':location.hash='';await p.onClose();break;

  case'plan':{if(fields.coverage==='area'&&(!String(fields.lat).trim()||!String(fields.lon).trim()||!String(fields.radius).trim()))throw Error('Enter a centre and radius for this area');p.setBusy('Checking map and directions downloads...');setPlan(null);const map=await p.request('/api/map-install/plan',{method:'POST',body:JSON.stringify({coverage:fields.coverage,depth:fields.depth,postcodes:fields.postcodes,lat:Number(fields.lat),lon:Number(fields.lon),radius:Number(fields.radius)})});let routing:any=null;if(fields.includeDirections&&['uk','area'].includes(fields.coverage))routing=await p.request('/api/directions/plan',{method:'POST',body:JSON.stringify({coverage:fields.coverage,lat:Number(fields.lat),lon:Number(fields.lon),radius:Number(fields.radius)})});if(routing&&map.requiredBytes+routing.requiredBytes>map.freeBytes)throw Error('Not enough free space for the combined map and directions. Choose a smaller area or less detail.');setPlan({...map,routing});break}
  case'install':await p.request('/api/map-install/start',{method:'POST',body:JSON.stringify({plan:plan.id})});if(plan.routing){queued.current={stage:'map',previousDir:directions?.installed?.id,choice:plan.routing.choice,bytes:plan.routing.downloadBytes,previous:status?.installed?.id};setQueueStage('Map first, then local directions. Keep this page open to continue both downloads.')}setPlan(null);setStatus(await p.request('/api/map-install'));break;
  case'cancel':queued.current=null;setQueueStage('');await p.request('/api/map-install/cancel',{method:'POST',body:'{}'});await p.request('/api/directions/cancel',{method:'POST',body:'{}'});break;
  case'connectRouter':case'disconnectRouter':await p.request('/api/map-install/routing',{method:'POST',body:JSON.stringify({url:a==='disconnectRouter'?'':fields.router.trim()})});setStatus(await p.request('/api/map-install'));p.setNotice(a==='disconnectRouter'?'Local directions disconnected':'Local driving directions connected');break;

  default:await p.common(a);

 }}catch(e:any){p.setError(e.message)}finally{p.setBusy('')}};

 const progress=status?.progress,installed=status?.installed||p.config.installed,working=!!progress?.busy||!!directions?.progress?.busy||!!queueStage;

 const data={directionsChoice:['uk','area'].includes(fields.coverage),queueStage,owner:!!p.config.canInstall,installed,installedLabel:installed?label(installed.choice):'',installedSize:installed?size(installed.bytes)+' on disk · '+installed.features.toLocaleString()+' searchable places':'',area:fields.coverage==='area',canPlan:!working,postcodeChoice:['uk','area','europe'].includes(fields.coverage),coverages:[{value:'world',label:'World overview'},{value:'uk',label:'United Kingdom'},{value:'europe',label:'Europe'},{value:'area',label:'Choose a local area'}],depths:[{value:'overview',label:'Overview · towns and main roads'},...(fields.coverage!=='world'?[{value:'region',label:'Regional · more roads and place names'}]:[]),...(fields.coverage==='area'?[{value:'streets',label:'Street detail · local roads and nearby places'}]:[])],depthText:fields.depth==='streets'?'Highest local detail. A smaller radius reduces storage and download time.':fields.depth==='region'?'More detail than overview, suitable for travelling across a region.':'Smallest download. Zooming in enlarges the overview; local POIs may be absent.',plan:!!plan&&!working,planLabel:plan?label(plan.choice):'',downloadLabel:plan?'Download: '+size(plan.downloadBytes+(plan.routing?.downloadBytes||0))+(plan.routing?' including local directions':'')+' · '+plan.tileCount.toLocaleString()+' tiles':'',spaceLabel:plan?'Free space required for safe installation: '+size(plan.requiredBytes+(plan.routing?.requiredBytes||0)):'',progress:progress&&progress.phase!=='choose',working,phaseLabel:working?'Installing your map':progress?.phase==='ready'?'Ready':progress?.phase==='cancelled'?'Cancelled':'Download needs attention',progressMessage:directions?.progress?.busy?directions.progress.message:progress?.message,progressBytes:progress?.total?size(progress.done)+' of '+size(progress.total)+' downloaded':'',cancelAllowed:working&&!!p.config.canInstall,routerStatus:status?.routing?'Advanced OSRM driving service connected.':'No advanced service connected.'};

 return p.render(React.createElement(Screen,{data,onAction:action}),{showHeader:true,context:{fields,setField:change,disabled:!!p.busy||working}},action,{'directions-download':React.createElement(DirectionsInstall,{request:p.request,owner:p.config.canInstall,status:directions,centre:p.centre,onReady:setDirections})});

}

