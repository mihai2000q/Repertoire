import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import IsAlreadyAuthenticated from './router/IsAlreadyAuthenticated'
import SignUp from './views/SignUp'
import RequireAuthentication from './router/RequireAuthentication'
import Home from './views/Home'
import NotFound from './views/NotFound'
import Unauthorized from './views/Unauthorized'
import { ReactElement } from 'react'
import Main from './views/Main'
import SignIn from './views/SignIn'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { Provider } from 'react-redux'
import { store } from './state/store'
import './index.css'
import '@mantine/core/styles.css'
import Songs from './views/Songs'

function App(): ReactElement {
  return (
    <div className={'app'}>
      <Provider store={store}>
        <MantineProvider theme={theme} forceColorScheme={'light'}>
          <BrowserRouter>
            <Routes>
              <Route path={'/'} element={<Navigate to={'home'} replace />} />

              <Route element={<IsAlreadyAuthenticated />}>
                <Route path={'sign-in'} element={<SignIn />} />
                <Route path={'sign-up'} element={<SignUp />} />
              </Route>

              <Route element={<RequireAuthentication />}>
                <Route element={<Main />}>
                  <Route path={'home'} element={<Home />} />
                  <Route path={'songs'} element={<Songs />} />

                  {/* Errors */}
                  <Route path={'401'} element={<Unauthorized />} />
                  <Route path={'404'} element={<NotFound />} />
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

// noinspection JSUnusedGlobalSymbols
export default App
