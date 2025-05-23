import './index.css'
import '@mantine/core/styles.css'
import '@mantine/dates/styles.css'
import { emotionTransform, MantineEmotionProvider } from '@mantine/emotion'
import { emotionCache } from './cache.ts'
import CustomizedToastContainer from './components/@ui/toast/CustomizedToastContainer.tsx'
import { ReactElement } from 'react'
import { BrowserRouter, HashRouter, Navigate, Route, Routes } from 'react-router-dom'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { Provider } from 'react-redux'
import { store } from './state/store'
import IsAlreadyAuthenticated from './router/IsAlreadyAuthenticated'
import SignUp from './views/SignUp'
import RequireAuthentication from './router/RequireAuthentication'
import Home from './views/Home'
import NotFound from './views/NotFound'
import Unauthorized from './views/Unauthorized'
import Main from './views/Main'
import SignIn from './views/SignIn'
import Songs from './views/Songs.tsx'
import Albums from './views/Albums.tsx'
import Artists from './views/Artists.tsx'
import Artist from './views/Artist.tsx'
import Album from './views/Album.tsx'
import Song from './views/Song.tsx'
import Playlists from './views/Playlists.tsx'
import Playlist from './views/Playlist.tsx'
import useIsDesktop from './hooks/useIsDesktop.ts'
import { MainScrollProvider } from './context/MainScrollContext.tsx'

function App(): ReactElement {
  const isDesktop = useIsDesktop()
  const Router = isDesktop ? HashRouter : BrowserRouter

  return (
    <div className={'app'}>
      <Provider store={store}>
        <MainScrollProvider>
          <MantineProvider
            theme={theme}
            forceColorScheme={'light'}
            stylesTransform={emotionTransform}
          >
            <MantineEmotionProvider cache={emotionCache}>
              <Router>
                <CustomizedToastContainer />
                <Routes>
                  <Route element={<Main />}>
                    <Route path={'/'} element={<Navigate to={'home'} replace />} />

                    <Route element={<IsAlreadyAuthenticated />}>
                      <Route path={'sign-in'} element={<SignIn />} />
                      <Route path={'sign-up'} element={<SignUp />} />
                    </Route>

                    <Route element={<RequireAuthentication />}>
                      <Route path={'home'} element={<Home />} />
                      <Route path={'artists'} element={<Artists />} />
                      <Route path={'artist/:id'} element={<Artist />} />
                      <Route path={'albums'} element={<Albums />} />
                      <Route path={'album/:id'} element={<Album />} />
                      <Route path={'songs'} element={<Songs />} />
                      <Route path={'song/:id'} element={<Song />} />
                      <Route path={'playlists'} element={<Playlists />} />
                      <Route path={'playlist/:id'} element={<Playlist />} />

                      {/* Errors */}
                      <Route path={'401'} element={<Unauthorized />} />
                      <Route path={'404'} element={<NotFound />} />
                      <Route path={'*'} element={<Navigate to={'404'} replace />} />
                    </Route>
                  </Route>
                </Routes>
              </Router>
            </MantineEmotionProvider>
          </MantineProvider>
        </MainScrollProvider>
      </Provider>
    </div>
  )
}

// noinspection JSUnusedGlobalSymbols
export default App
