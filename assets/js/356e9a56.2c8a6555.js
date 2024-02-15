"use strict";(self.webpackChunksite=self.webpackChunksite||[]).push([[8656],{7964:(e,t,o)=>{o.r(t),o.d(t,{assets:()=>l,contentTitle:()=>r,default:()=>u,frontMatter:()=>i,metadata:()=>a,toc:()=>d});var n=o(2488),s=o(2780);const i={sidebar_position:3,id:"settlement"},r="Rollup Settlement",a={id:"background/settlement",title:"Rollup Settlement",description:'Rollups come in two main flavors: optimistic rollups and ZK rollups. They have different "settlement" mechanisms for finalizing the state of their rollup.',source:"@site/versioned_docs/version-0.1.0/background/settlement.md",sourceDirName:"background",slug:"/background/settlement",permalink:"/docs/background/settlement",draft:!1,unlisted:!1,editUrl:"https://github.com/omni-network/omni/docs/versioned_docs/version-0.1.0/background/settlement.md",tags:[],version:"0.1.0",sidebarPosition:3,frontMatter:{sidebar_position:3,id:"settlement"},sidebar:"oldSidebar",previous:{title:"Ethereum and Rollups",permalink:"/docs/background/rollups"},next:{title:"Contracts Overview",permalink:"/docs/developers/contracts"}},l={},d=[];function c(e){const t={h1:"h1",p:"p",...(0,s.M)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.h1,{id:"rollup-settlement",children:"Rollup Settlement"}),"\n",(0,n.jsx)(t.p,{children:'Rollups come in two main flavors: optimistic rollups and ZK rollups. They have different "settlement" mechanisms for finalizing the state of their rollup.'}),"\n",(0,n.jsx)(t.p,{children:"Optimistic rollups operate with fault proofs. A privileged actor can submit state updates to the rollup's settlement contract on Ethereum Layer 1. Over a period of ~7 days, anyone can dispute the validity of this state update by submitting a proof that the update is invalid. As long as there is 1 honest actor watching this rollup, invalid state updates will be challenged and thrown out by the end of the 7 day window."}),"\n",(0,n.jsx)(t.p,{children:'For ZK (zero knowledge) rollups, settlement works differently. Instead of the 7 day challenge window, a "prover" runs a sophisticated computation to generate a validity proof. This proof is submitted along with the state updates and provides a mathematical guarantee that the state update is valid.'}),"\n",(0,n.jsx)(t.p,{children:"Both systems come with trade-offs. The 7-day optimistic window introduces friction for users who wish to withdraw their assets from the rollup back to Layer 1. ZK rollups rely on complex math and heavy computation which will make decentralizing them over time more difficult."}),"\n",(0,n.jsx)(t.p,{children:"Omni introduces a unified settlement system across both rollup categories and decreases settlement time to establish a faster interoperability standard. You can read on in the Protocol section for more details."})]})}function u(e={}){const{wrapper:t}={...(0,s.M)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(c,{...e})}):c(e)}},2780:(e,t,o)=>{o.d(t,{I:()=>a,M:()=>r});var n=o(6651);const s={},i=n.createContext(s);function r(e){const t=n.useContext(i);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:r(e.components),n.createElement(i.Provider,{value:t},e.children)}}}]);