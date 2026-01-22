import React from 'react'
import ReactDOM from 'react-dom/client'
import logo from './assets/logo.ico'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '97vh',
        width: '97vw',
        aspectRatio: 1
      }}
    >
      <img src={logo} alt={'logo'} />
    </div>
  </React.StrictMode>
)
