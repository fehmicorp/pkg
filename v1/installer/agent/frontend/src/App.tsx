import { TitleBar } from "./components/TitleBar"

function App() {
  return (
    <div className='flex h-screen w-screen items-center justify-center bg-slate-950 p-4 font-sans text-slate-100 select-none'>
      <div className="flex h-[560px] w-[820px] flex-col overflow-hidden rounded-xl border border-slate-800 bg-slate-900 shadow-2xl shadow-cyan-950/20">
        <TitleBar/>
      </div>
    </div>
  )
}

export default App
