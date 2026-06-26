import './index.css'
import '@mantine/core/styles.css'
import '@mantine/dates/styles.css'
import { emotionTransform, MantineEmotionProvider } from '@mantine/emotion'
import { emotionCache } from './cache.ts'
import CustomizedToastContainer from './components/@ui/toast/CustomizedToastContainer.tsx'
import { ReactElement } from 'react'
import { BrowserRouter, HashRouter, Navigate, Route, Routes } from 'react-router-dom'
import { MantineProvider, v8CssVariablesResolver } from '@mantine/core'
import { theme } from './theme/theme'
import { Provider } from 'react-redux'
import { store } from './state/store'
import IsAlreadyAuthenticated from './router/IsAlreadyAuthenticated'
import SignUp from './features/auth/features/sign-up/SignUp.tsx'
import RequireAuthentication from './router/RequireAuthentication'
import Home from './features/home/Home.tsx'
import NotFound from './features/not-found/NotFound.tsx'
import Unauthorized from './features/unauthorized/Unauthorized.tsx'
import Main from './features/main/Main.tsx'
import SignIn from './features/auth/features/sign-in/SignIn.tsx'
import Songs from './features/songs/Songs.tsx'
import Albums from './features/albums/Albums.tsx'
import Artists from './features/artists/Artists.tsx'
import Artist from './features/artists/features/artist/Artist.tsx'
import Album from './features/albums/features/album/Album.tsx'
import Song from './features/songs/features/song/Song.tsx'
import Playlists from './features/playlists/Playlists.tsx'
import Playlist from './features/playlists/features/playlist/Playlist.tsx'
import useIsDesktop from './hooks/useIsDesktop.ts'
import dayjs from 'dayjs'
import isToday from 'dayjs/plugin/isToday'
import isYesterday from 'dayjs/plugin/isYesterday'
import 'dayjs/locale/en-gb'

dayjs.extend(isToday)
dayjs.extend(isYesterday)
dayjs.locale('en-gb')

function App(): ReactElement {
  const isDesktop = useIsDesktop()
  const Router = isDesktop ? HashRouter : BrowserRouter

  return (
    <div className={'app'}>
      <Provider store={store}>
        <MantineProvider
          theme={theme}
          forceColorScheme={'light'}
          stylesTransform={emotionTransform}
          cssVariablesResolver={v8CssVariablesResolver} // TODO: Remove and adapt styling
          deduplicateInlineStyles
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
      </Provider>
    </div>
  )
}

// noinspection JSUnusedGlobalSymbols
export default App
