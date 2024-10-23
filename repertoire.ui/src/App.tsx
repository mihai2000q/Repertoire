import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import IsAlreadyAuthenticated from './router/IsAlreadyAuthenticated'
import RequireAuthentication from './router/RequireAuthentication'
import { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import './index.css'
import '@mantine/core/styles.css'
import Main from './views/Main'
import SignIn from './views/SignIn'
import SignUp from './views/SignUp'
import Songs from './views/Songs'
import Home from './views/Home'
import NotFound from './views/NotFound'
import Unauthorized from './views/Unauthorized'

const queryClient = new QueryClient()

function App(): ReactElement {
  return (
    <div className={'app'}>
      <QueryClientProvider client={queryClient}>
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
      </QueryClientProvider>
    </div>
  )
}

// noinspection JSUnusedGlobalSymbols
export default App
