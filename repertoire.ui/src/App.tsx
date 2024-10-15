import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import IsAlreadyAuthenticated from './router/IsAlreadyAuthenticated'
import SignUpView from './views/SignUpView'
import RequireAuthentication from './router/RequireAuthentication'
import HomeView from './views/HomeView'
import NotFoundView from './views/NotFoundView'
import UnauthorizedView from './views/UnauthorizedView'
import { ReactElement } from 'react'
import MainView from './views/MainView'
import SignInView from './views/SignInView'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { Provider } from 'react-redux'
import { store } from './state/store'
import './index.css'

function App(): ReactElement {
  return (
    <div className={'app'}>
      <Provider store={store}>
        <MantineProvider theme={theme}>
          <BrowserRouter>
            <Routes>
              <Route path={'/'} element={<Navigate to={'home'} replace />} />

              <Route element={<IsAlreadyAuthenticated />}>
                <Route path={'sign-in'} element={<SignInView />} />
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
        </MantineProvider>
      </Provider>
    </div>
  )
}

export default App
