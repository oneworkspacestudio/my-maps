import React,{useEffect,useState} from 'react';

import {MapInstall as Screen} from './generated/MapInstall';

const size=(n:number)=>n>=1073741824?(n/1073741824).toFixed(2)+' GB':(n/1048576).toFixed(1)+' MB';

const coverageName=(c:any)=>({world:'World',uk:'United Kingdom',europe:'Europe',area:'Local area'}[c.coverage as string]||c.coverage);

const label=(c:any)=>coverageName(c)+' · '+({overview:'Overview',region:'Regional detail',streets:'Street detail'}[c.depth as string]||c.depth)+(c.coverage==='area'?' · '+c.radius+' km radius':'');

export function MapInstall(p:any){

 const [fields,setFields]=useState<any>({coverage:p.config.installed?.choice.coverage||'uk',depth:p.config.installed?.choice.depth||'overview',lat:Number(p.centre?.lat??51.505).toFixed(5),lon:Number(p.centre?.lon??-.09).toFixed(5),radius:String(p.config.installed?.choice.radius||10),postcodes:!!p.config.installed?.choice.postcodes,router:''}),[status,setStatus]=useState<any>(null),[plan,setPlan]=useState<any>(null);

 useEffect(()=>{let live=true;const load=()=>p.request('/api/map-install').then((s:any)=>{if(live)setStatus(s)}).catch((e:any)=>{if(live)p.setError(e.message)});load();const timer=setInterval(load,1500);return()=>{live=false;clearInterval(timer)}},[]);

 const change=(key:string,value:any)=>{setPlan(null);setFields((f:any)=>{const next={...f,[key]:value};if(key==='coverage'&&value==='world'){next.depth='overview';next.postcodes=false}if(key==='coverage'&&value!=='area'&&next.depth==='streets')next.depth='region';return next})};

 const action=async(a:string)=>{p.setError('');try{switch(a){

  case'mapBack':case'openMap':location.hash='';await p.onClose();break;

  case'plan':if(fields.coverage==='area'&&(!String(fields.lat).trim()||!String(fields.lon).trim()||!String(fields.radius).trim()))throw Error('Enter a centre and radius for this area');p.setBusy('Checking the selected map download...');setPlan(await p.request('/api/map-install/plan',{method:'POST',body:JSON.stringify({coverage:fields.coverage,depth:fields.depth,postcodes:fields.postcodes,lat:Number(fields.lat),lon:Number(fields.lon),radius:Number(fields.radius)})}));break;

  case'install':await p.request('/api/map-install/start',{method:'POST',body:JSON.stringify({plan:plan.id})});setPlan(null);setStatus(await p.request('/api/map-install'));break;

  case'cancel':await p.request('/api/map-install/cancel',{method:'POST',body:'{}'});break;

  case'connectRouter':case'disconnectRouter':await p.request('/api/map-install/routing',{method:'POST',body:JSON.stringify({url:a==='disconnectRouter'?'':fields.router.trim()})});setStatus(await p.request('/api/map-install'));p.setNotice(a==='disconnectRouter'?'Local directions disconnected':'Local driving directions connected');break;

  default:await p.common(a);

 }}catch(e:any){p.setError(e.message)}finally{p.setBusy('')}};

 const progress=status?.progress,installed=status?.installed||p.config.installed,working=!!progress?.busy;

 const data={owner:!!p.config.canInstall,installed,installedLabel:installed?label(installed.choice):'',installedSize:installed?size(installed.bytes)+' on disk · '+installed.features.toLocaleString()+' searchable places':'',area:fields.coverage==='area',canPlan:!working,postcodeChoice:['uk','area','europe'].includes(fields.coverage),coverages:[{value:'world',label:'World overview'},{value:'uk',label:'United Kingdom'},{value:'europe',label:'Europe'},{value:'area',label:'Choose a local area'}],depths:[{value:'overview',label:'Overview · towns and main roads'},...(fields.coverage!=='world'?[{value:'region',label:'Regional · more roads and place names'}]:[]),...(fields.coverage==='area'?[{value:'streets',label:'Street detail · local roads and nearby places'}]:[])],depthText:fields.depth==='streets'?'Highest local detail. A smaller radius reduces storage and download time.':fields.depth==='region'?'More detail than overview, suitable for travelling across a region.':'Smallest download. Zooming in enlarges the overview; local POIs may be absent.',plan:!!plan&&!working,planLabel:plan?label(plan.choice):'',downloadLabel:plan?'Download: '+size(plan.downloadBytes)+' · '+plan.tileCount.toLocaleString()+' tiles':'',spaceLabel:plan?'Free space required for safe installation: '+size(plan.requiredBytes):'',progress:progress&&progress.phase!=='choose',working,phaseLabel:working?'Installing your map':progress?.phase==='ready'?'Ready':progress?.phase==='cancelled'?'Cancelled':'Download needs attention',progressMessage:progress?.message,progressBytes:progress?.total?size(progress.done)+' of '+size(progress.total)+' downloaded':'',cancelAllowed:working&&!!p.config.canInstall,routerStatus:status?.routing?'Local driving service connected.':'Directions are hidden until a local routing service is connected.'};

 return p.render(React.createElement(Screen,{data,onAction:action}),{showHeader:true,context:{fields,setField:change,disabled:!!p.busy||working}},action);

}

