import{_ as g,r as s,V as k,U as w,f as y,c as $,o as x,h as C,w as E,i as u,t as p}from"./index-SiRIT6sV.js";const B={class:"app-container"},I={class:"prompt"},b=`Last login: Fri Nov  1 15:53:18 on ttys008
You have new mail.`,S={__name:"monitoring",setup(T){const n=s(""),t=s(""),r=s("phearunum@Picha-Mac-mini ~ % "),l=s([]);let a=s(-1);t.value=`${b}
${r.value}`;let o=null;const v=()=>{o=new WebSocket("ws://localhost:1212/terminal"),o.onopen=()=>{t.value+=`
[Connected to server]
`},o.onmessage=e=>{t.value+=`
${e.data}`,t.value+=`
${r.value}`,i()},o.onerror=e=>{console.error("WebSocket error:",e),t.value+=`
[Error connecting to server]`},o.onclose=()=>{t.value+=`
[Disconnected from server]`}},m=e=>{n.value=e.target.innerText},d=e=>{if(e.key==="ArrowUp"){if(a.value<l.value.length-1){a.value++,n.value=l.value[l.value.length-1-a.value]||"",c();return}}else if(e.key==="ArrowDown"&&a.value>0){a.value--,n.value=l.value[l.value.length-1-a.value]||"",c();return}e.key==="Enter"&&f()},f=()=>{n.value.trim()!==""&&(t.value+=`
${r.value}${n.value}`,o.send(n.value),l.value.push(n.value),a.value=0,n.value="",c(),i())},c=()=>{const e=$refs.input;e.innerText="",e.focus()},i=()=>{const e=$refs.terminalOutput;e.scrollTop=e.scrollHeight},_=()=>{$refs.input.focus()};return k(()=>{o&&o.close()}),w(()=>{v()}),(e,D)=>{const h=y("el-card");return x(),$("div",B,[C(h,{class:"terminal-container"},{default:E(()=>[u("div",{class:"terminal-output",ref:"terminalOutput",onClick:_},[u("pre",null,p(t.value),1),u("span",I,p(r.value),1),u("span",{contenteditable:"true",class:"input",ref:"input",onInput:m,onKeydown:d},null,544)],512)]),_:1})])}}},O=g(S,[["__scopeId","data-v-3ce04545"]]);export{O as default};
