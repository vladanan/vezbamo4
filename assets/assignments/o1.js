var bundle=(()=>{function u(e,a){return Math.floor(Math.random()*(a-e+1))+e}function C(e){let a="neko pitanje",l=e.novo_pitanje.a,n=e.novo_pitanje.op2,m=e.novo_pitanje.b;return a=React.createElement(React.Fragment,null,React.createElement("img",{className:"mb-3",width:"525",src:"/static/assignments/brojevna-prava2.png"}),React.createElement("label",{className:"text-2xl ml-2",htmlFor:"r"},l," ",n," ",m," = "),React.createElement("input",{onChange:d=>e.setOdgovor(d.target.value),className:"text-2xl w-16  bg-lime-300 px-1 border border-lime-700 rounded-md",type:"number",min:"0",max:"10",id:"r",name:"r",defaultValue:0}),React.createElement("span",{className:"text-2xl ml-2 font-bold text-green-700 "},"\xA0\xA0",e.resenje)),React.createElement(React.Fragment,null,React.createElement("div",{className:"p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl"},a),React.createElement("div",{hidden:!e.vidljiv_odgovor,className:"p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl"},e.rezultat))}function O(e){let a="neko pitanje",l=e.novo_pitanje.a,n=e.novo_pitanje.b,m=e.novo_pitanje.t11,b=e.novo_pitanje.t22,d=e.novo_pitanje.t33,c,_;return l>n?(c=l-n,_=m+" \u0438\u043C\u0430 "+l+" "+d+" a "+b+" "+c+" \u043C\u0430\u045A\u0435. \u041A\u043E\u043B\u0438\u043A\u043E \u0438\u043C\u0430\u0458\u0443 \u0437\u0430\u0458\u0435\u0434\u043D\u043E?"):(c=n-l,_=m+" \u0438\u043C\u0430 "+l+" "+d+" a "+b+" "+c+" \u0432\u0438\u0448\u0435. \u041A\u043E\u043B\u0438\u043A\u043E \u0438\u043C\u0430\u0458\u0443 \u0443\u043A\u0443\u043F\u043D\u043E?"),a=React.createElement(React.Fragment,null,React.createElement("img",{className:"mb-3",width:"525",src:"/static/assignments/brojevna-prava2.png"}),React.createElement("label",{className:"text-2xl",htmlFor:"r"},_," "),React.createElement("input",{onChange:g=>e.setOdgovor(g.target.value),className:"text-2xl w-16 bg-lime-300 px-1 border border-lime-700 rounded-md",type:"number",min:"0",max:"10",id:"r",name:"r",defaultValue:0}),React.createElement("span",{className:"text-2xl ml-2 font-bold text-green-700 "},"\xA0\xA0",e.resenje)),React.createElement(React.Fragment,null,React.createElement("div",{className:"p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl"},a),React.createElement("div",{hidden:!e.vidljiv_odgovor,className:"p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl"},e.rezultat))}function A(e){let[a,l]=React.useState("o1m_1_10"),[n,m]=React.useState(!1),[b,d]=React.useState({a:2,op2:"+",b:2,t11:"\u041C\u0430\u0440\u043A\u043E",t22:"\u041C\u0430\u0448\u0430",t33:"\u043A\u0430\u043C\u0438\u043E\u043D\u0447\u0438\u045B\u0430",kombi:"\u041A\u043B\u0438\u043A\u043D\u0438 \u043D\u0430: \u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A!"}),[c,_]=React.useState(null),[h,g]=React.useState(null),[p,y]=React.useState(null),x=t=>{l(t),t=="o1m_1_10"?d({a:0,op2:"+",b:0,t11:"\u041C\u0438\u043B\u043E\u0448",t22:"\u0408\u0435\u043B\u0435\u043D\u0430",t33:"\u043A\u043E\u043B\u0430\u0447\u0430",kombi:"\u041A\u043B\u0438\u043A\u043D\u0438 \u043D\u0430: \u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A!"}):d({a:0,op2:"*",b:0,t11:"\u041C\u0438\u043B\u043E\u0448",t22:"\u0408\u0435\u043B\u0435\u043D\u0430",t33:"\u043A\u043E\u043B\u0430\u0447\u0430",kombi:"\u041A\u043B\u0438\u043A\u043D\u0438 \u043D\u0430: \u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A!"})},k=()=>{let t;switch(a){case"o1m_1_10":var{a:s,op2:f,b:o}=b,v=c,r="",i=0;if(f=="-")var r=s-o;else var r=s+o;r==v?i=1:i=0,g(r),i==1?t=React.createElement(React.Fragment,null,React.createElement("p",{style:{textAlign:"center",background:"skyblue",fontSize:"30px"}},"\u2714")):t=React.createElement(React.Fragment,null,React.createElement("p",{style:{background:"red",textAlign:"center",fontSize:"30px"}},"\u2718"));break;case"o1m_1_10txt":var{a:s,b:o}=b,v=c,r=s+o,i=0;r==v?i=1:i=0,g(r),i==1?t=React.createElement(React.Fragment,null,React.createElement("p",{style:{textAlign:"center",background:"skyblue",fontSize:"30px"}},"\u2714")):t=React.createElement(React.Fragment,null,React.createElement("p",{style:{background:"red",textAlign:"center",fontSize:"30px"}},"\u2718"));break;default:}n==!0?m(!1):(m(!0),y(t))},j=()=>{switch(a){case"o1m_1_10":let t=Math.floor(Math.random()*11),f=Math.floor(Math.random()*2),s="",o=Math.floor(Math.random()*11);if(f==0){if(t<o){let S=t;t=o,o=S}}else for(;t+o>10;)t=Math.floor(Math.random()*11),o=Math.floor(Math.random()*11);f==0?s="-":s="+",g(null),d({a:t,op2:s,b:o,t11:"\u041C\u0438\u043B\u043E\u0448",t22:"\u0408\u0435\u043B\u0435\u043D\u0430",t33:"\u043A\u043E\u043B\u0430\u0447\u0430",kombi:"\u041A\u043B\u0438\u043A\u043D\u0438 \u043D\u0430: \u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A!"});break;case"o1m_1_10txt":let v=u(2,10),r=u(2,10),i=["\u041C\u0430\u0440\u043A\u043E","M\u0430\u0458\u043C\u0443\u043D","\u041C\u0438\u043B\u043E\u0448","K\u0440\u043E\u043A\u043E\u0434\u0438\u043B","\u0421\u0442\u0435\u0444\u0430\u043D","\u0421\u043B\u0435\u043F\u0438 \u043C\u0438\u0448","\u0418\u043B\u0438\u0458\u0430","\u0421\u043B\u043E\u043D","\u0408\u043E\u0432\u0430\u043D","\u0413\u043E\u0440\u0438\u043B\u0430"],N=["\u041C\u0430\u0448\u0430","\u0436\u0438\u0440\u0430\u0444\u0430","\u0408\u0435\u043B\u0435\u043D\u0430","\u0445\u0438\u0458\u0435\u043D\u0430","\u0418\u0432\u0430\u043D\u0430","\u0447\u0430\u043F\u0459\u0430","\u041C\u0438\u043B\u0435\u043D\u0430","\u043A\u043E\u043A\u043E\u0448\u043A\u0430","\u041D\u0430\u0442\u0430\u0448\u0430","\u043E\u0432\u0446\u0430"],w=["\u043A\u0430\u043C\u0438\u043E\u043D\u0447\u0438\u045B\u0430","\u0435\u043A\u0441\u0435\u0440\u0430","\u043A\u043E\u043B\u0430\u0447\u0430","\u0447\u0435\u043A\u0438\u045B\u0430","\u0434\u0438\u043D\u0430\u0440\u0430","\u0444\u043B\u043E\u043C\u0430\u0441\u0442\u0435\u0440\u0430","\u043B\u0438\u0437\u0430\u043B\u0430","\u0446\u0432\u0435\u0442\u0438\u045B\u0430","\u043B\u0435\u043F\u0442\u0438\u0440\u0438\u045B\u0430","\u045B\u0435\u0432\u0430\u043F\u0447\u0438\u045B\u0430"],z=i[u(0,9)],M=N[u(0,9)],R=w[u(0,9)];for(;v+r>10;)v=u(2,10),r=u(2,10);g(null),d({a:v,op2:"+",b:r,t11:z,t22:M,t33:R,kombi:"\u041A\u043B\u0438\u043A\u043D\u0438 \u043D\u0430: \u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A!"});break;default:break}m(!1)};return React.createElement(React.Fragment,null,a=="o1m_1_10"&&React.createElement(React.Fragment,null,React.createElement("div",null,e.id,e.name)),React.createElement("p",{className:"text-2xl mt-5 ml-2"},"\u041E\u0434\u0430\u0431\u0435\u0440\u0438 \u0437\u0430\u0434\u0430\u0442\u043A\u0435:"),React.createElement("div",{className:"mx-2 p-2 text-black max-w-sm border-2 border-sky-700 rounded-md bg-gradient-to-r from-blue-50 to-transparent"},React.createElement("label",{className:a=="o1m_1_10"?"relative p-1 border-2 border-black rounded-md bg-sky-400":"relative p-1 border-2 border-black rounded-md bg-gray-400",htmlFor:"o1m_1_10"},"\u0421\u0430\u0431\u0438\u0440\u0430\u045A\u0435 \u0438 \u043E\u0434\u0443\u0437\u0438\u043C\u0430\u045A\u0435 \u043E\u0434 1 \u0434\u043E 10"),React.createElement("input",{onClick:t=>x("o1m_1_10"),className:"relative ml-2 mt-2",type:"radio",id:"o1m_1_10",name:"radio_z",value:"o1m_1_10",defaultChecked:!0}),React.createElement("br",null),React.createElement("label",{className:a=="o1m_1_10txt"?"relative p-1 border-2 border-black rounded-md bg-gradient-to-br from-red-500 via-yellow-100 to-green-300":"relative p-1 border-2 border-black rounded-md bg-gray-400",htmlFor:"o1m_1_10txt"},"\u0421\u043C\u0435\u0448\u043D\u043E \u0441\u0430\u0431. \u0438 \u043E\u0434\u0443\u0437. \u043E\u0434 1 \u0434\u043E 10"),React.createElement("input",{onClick:t=>x("o1m_1_10txt"),className:"relative ml-2 mt-5",type:"radio",id:"o1m_1_10txt",name:"radio_z",value:"o1m_1_10txt"}),React.createElement("br",null),a=="o1m_1_10"&&React.createElement(React.Fragment,null,React.createElement("br",null),React.createElement("button",{onClick:k,className:`ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200`},"\u0420\u0435\u0448\u0435\u045A\u0435"),React.createElement("button",{onClick:j,className:`ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400`},"\u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A"),React.createElement("br",null),React.createElement(C,{novo_pitanje:b,vidljiv_odgovor:n,setOdgovor:_,rezultat:p,resenje:h})),a=="o1m_1_10txt"&&React.createElement(React.Fragment,null,React.createElement("br",null),React.createElement("button",{onClick:k,className:`ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200`},"\u0420\u0435\u0448\u0435\u045A\u0435"),React.createElement("button",{onClick:j,className:`ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400`},"\u041D\u043E\u0432\u0438 \u0437\u0430\u0434\u0430\u0442\u0430\u043A"),React.createElement("br",null),React.createElement(O,{novo_pitanje:b,vidljiv_odgovor:n,setOdgovor:_,rezultat:p,resenje:h}))),React.createElement("br",null))}var F=document.getElementById("root"),V=ReactDOM.createRoot(F);V.render(React.createElement(A,null));})();
