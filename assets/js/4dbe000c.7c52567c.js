"use strict";(self.webpackChunk=self.webpackChunk||[]).push([[489],{5318:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>m});var r=n(7378);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,i=function(e,t){if(null==e)return{};var n,r,i={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var s=r.createContext({}),c=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},d=function(e){var t=c(e.components);return r.createElement(s.Provider,{value:t},e.children)},u="mdxType",p={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},f=r.forwardRef((function(e,t){var n=e.components,i=e.mdxType,o=e.originalType,s=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),u=c(n),f=i,m=u["".concat(s,".").concat(f)]||u[f]||p[f]||o;return n?r.createElement(m,a(a({ref:t},d),{},{components:n})):r.createElement(m,a({ref:t},d))}));function m(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=n.length,a=new Array(o);a[0]=f;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l[u]="string"==typeof e?e:i,a[1]=l;for(var c=2;c<o;c++)a[c]=n[c];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}f.displayName="MDXCreateElement"},5844:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>s,contentTitle:()=>a,default:()=>p,frontMatter:()=>o,metadata:()=>l,toc:()=>c});var r=n(5773),i=(n(7378),n(5318));const o={id:"linux",title:"KNDP for Linux",sidebar_posiion:2},a='<div align="center">   Install KNDP  </div>',l={unversionedId:"getting_started/install/linux",id:"getting_started/install/linux",title:"KNDP for Linux",description:"Prerequisites",source:"@site/docs/getting_started/install/linux.md",sourceDirName:"getting_started/install",slug:"/getting_started/install/linux",permalink:"/docs/getting_started/install/linux",draft:!1,tags:[],version:"current",frontMatter:{id:"linux",title:"KNDP for Linux",sidebar_posiion:2},sidebar:"myAutogeneratedSidebar",previous:{title:"Installation",permalink:"/docs/category/installation"},next:{title:"Development View",permalink:"/docs/development-view"}},s={},c=[{value:"Prerequisites",id:"prerequisites",level:3},{value:"Run with:",id:"run-with",level:3},{value:"!! DISCLAIMER !!",id:"-disclaimer-",level:3},{value:"If Docker was installed by the script and was not present before, it is recommended to log out and log back in to apply changes and run Docker without <code>sudo</code>.",id:"if-docker-was-installed-by-the-script-and-was-not-present-before-it-is-recommended-to-log-out-and-log-back-in-to-apply-changes-and-run-docker-without-sudo",level:4}],d={toc:c},u="wrapper";function p(e){let{components:t,...n}=e;return(0,i.kt)(u,(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"---install-kndp--"},(0,i.kt)("div",{align:"center"},"   Install KNDP  ")),(0,i.kt)("h3",{id:"prerequisites"},"Prerequisites"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"Linux-based operating system. This script is designed for Linux-based operating systems, specifically tested on Ubuntu."),(0,i.kt)("li",{parentName:"ol"},"Internet access for package downloads."),(0,i.kt)("li",{parentName:"ol"},"Administrative privileges. Ensure you have administrative privileges to execute certain commands (e.g., sudo).")),(0,i.kt)("h3",{id:"run-with"},"Run with:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"}," curl -sSf https://raw.githubusercontent.com/web-seven/kndp/release/0.1/scripts/install.sh | bash\n")),(0,i.kt)("h3",{id:"-disclaimer-"},"!! DISCLAIMER !!"),(0,i.kt)("h4",{id:"if-docker-was-installed-by-the-script-and-was-not-present-before-it-is-recommended-to-log-out-and-log-back-in-to-apply-changes-and-run-docker-without-sudo"},"If Docker was installed by the script and was not present before, it is recommended to log out and log back in to apply changes and run Docker without ",(0,i.kt)("inlineCode",{parentName:"h4"},"sudo"),"."))}p.isMDXComponent=!0}}]);