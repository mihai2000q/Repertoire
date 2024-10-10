import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import IsAlreadyAuthenticated from '@renderer/router/IsAlreadyAuthenticated'
import SignUpView from '@renderer/views/SignUpView'
import RequireAuthentication from '@renderer/router/RequireAuthentication'
import HomeView from '@renderer/views/HomeView'
import NotFoundView from '@renderer/views/NotFoundView'
import UnauthorizedView from '@renderer/views/UnauthorizedView'
import { ReactElement } from 'react'
import MainView from '@renderer/views/MainView'

function App(): ReactElement {
  return (
    <div className={'app'}>
      <BrowserRouter>
        <Routes>
          <Route path={'/'} element={<Navigate to={'home'} replace />} />

          <Route element={<IsAlreadyAuthenticated />}>
            <Route path={'sign-up'} element={<SignUpView />} />
          </Route>

          <Route element={<RequireAuthentication />}>
            <Route element={<MainView />}>
              <Route path={'home'} element={<HomeView />} />

              {/* Errors */}
              <Route path={'401'} element={<UnauthorizedView />} />
              <Route path={'404'} element={<NotFoundView />} />
              <Route path={'*'} element={<Navigate to={'404'} replace />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App
