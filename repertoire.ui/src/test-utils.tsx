import { render, renderHook, RenderHookResult, RenderResult } from '@testing-library/react'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { ReactNode } from 'react'
import { Provider } from 'react-redux'
import { EnhancedStore } from '@reduxjs/toolkit'
import { RootState, setupStore } from './state/store'
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router-dom'
import { emotionTransform, MantineEmotionProvider } from '@mantine/emotion'
import { ToastContainer } from 'react-toastify'
import Album from './types/models/Album.ts'
import Song, { SongSection, SongSettings } from './types/models/Song.ts'
import Artist from './types/models/Artist.ts'
import Order from './types/Order.ts'
import User from './types/models/User.ts'
import Playlist from './types/models/Playlist.ts'
import { AlbumFiltersMetadata, ArtistFiltersMetadata, PlaylistFiltersMetadata } from './types/models/FiltersMetadata.ts'

// Custom Matchers

// noinspection JSUnusedGlobalSymbols
expect.extend({
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  async toBeFormDataImage(received: string | Blob, expected: File) {
    if (!received)
      return {
        pass: false,
        message: () => `received form data image is not defined`
      }
    const { isNot } = this

    const capturedImage = received as File
    const equalType = capturedImage.type === expected.type
    const equalName = capturedImage.name === expected.name

    const content = await capturedImage.arrayBuffer()
    const expectedContent = await expected.arrayBuffer()
    let equalContent = true
    for (let i = 0; i < content.byteLength; i++) {
      if (content[i] !== expectedContent[i]) {
        equalContent = false
        break
      }
    }

    return {
      pass: equalType && equalName && equalContent,
      message: () => `received form data image is${isNot ? ' not' : ''} the expected image`
    }
  }
})

// Custom Renders

const MantineProviderComponent = ({ children }: { children: ReactNode }) => (
  <MantineProvider theme={theme} stylesTransform={emotionTransform}>
    <MantineEmotionProvider>{children}</MantineEmotionProvider>
  </MantineProvider>
)

export function withToastify(ui: ReactNode) {
  return (
    <>
      <ToastContainer />
      {ui}
    </>
  )
}

export function mantineRender(ui: ReactNode) {
  return render(ui, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <MantineProviderComponent>{children}</MantineProviderComponent>
    )
  })
}

export function routerRender(ui: ReactNode) {
  return render(ui, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <BrowserRouter>
        <MantineProviderComponent>{children}</MantineProviderComponent>
      </BrowserRouter>
    )
  })
}

export function reduxRender(
  ui: ReactNode,
  preloadedState?: Partial<RootState>
): [RenderResult, EnhancedStore] {
  const store = setupStore(preloadedState)

  return [
    render(ui, {
      wrapper: ({ children }: { children: ReactNode }) => (
        <Provider store={store}>
          <MantineProviderComponent>{children}</MantineProviderComponent>
        </Provider>
      )
    }),
    store
  ]
}

export function reduxRouterRender(
  ui: ReactNode,
  preloadedState?: Partial<RootState>
): [RenderResult, EnhancedStore] {
  const store = setupStore(preloadedState)

  return [
    render(ui, {
      wrapper: ({ children }: { children: ReactNode }) => (
        <Provider store={store}>
          <BrowserRouter>
            <MantineProviderComponent>{children}</MantineProviderComponent>
          </BrowserRouter>
        </Provider>
      )
    }),
    store
  ]
}

export function reduxMemoryRouterRender(
  ui: ReactNode,
  path: string,
  initialEntries: string[],
  preloadedState?: Partial<RootState>
): [RenderResult, EnhancedStore] {
  const store = setupStore(preloadedState)

  return [
    render(ui, {
      wrapper: ({ children }: { children: ReactNode }) => (
        <Provider store={store}>
          <MantineProviderComponent>
            <MemoryRouter initialEntries={initialEntries}>
              <Routes>
                <Route path={path} element={children} />
              </Routes>
            </MemoryRouter>
          </MantineProviderComponent>
        </Provider>
      )
    }),
    store
  ]
}

// Hooks

export function mantineRenderHook<T>(hook: (props: T) => T) {
  return renderHook(hook, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <MantineProviderComponent>{children}</MantineProviderComponent>
    )
  })
}

export function routerRenderHook<T>(hook: (props: T) => T): RenderHookResult<T, T> {
  return renderHook(hook, {
    wrapper: ({ children }: { children: ReactNode }) => <BrowserRouter>{children}</BrowserRouter>
  })
}

export function reduxRenderHook<T>(
  hook: (props: T) => T,
  preloadedState?: Partial<RootState>
): [RenderHookResult<T, T>, EnhancedStore] {
  const store = setupStore(preloadedState)

  return [
    renderHook(hook, {
      wrapper: ({ children }: { children: ReactNode }) => (
        <Provider store={store}>{children}</Provider>
      )
    }),
    store
  ]
}

export function reduxRouterRenderHook<T>(
  hook: (props: T) => T,
  preloadedState?: Partial<RootState>
): [RenderHookResult<T, T>, EnhancedStore] {
  const store = setupStore(preloadedState)

  return [
    renderHook(hook, {
      wrapper: ({ children }: { children: ReactNode }) => (
        <Provider store={store}>
          <BrowserRouter>{children}</BrowserRouter>
        </Provider>
      )
    }),
    store
  ]
}

// Empty Types

export const emptyUser: User = {
  createdAt: '',
  email: '',
  id: '',
  name: '',
  updatedAt: ''
}

export const emptyArtist: Artist = {
  id: '',
  name: '',
  isBand: false,
  createdAt: '',
  updatedAt: '',
  albums: [],
  songs: [],
  bandMembers: []
}

export const emptyAlbum: Album = {
  createdAt: '',
  id: '',
  songs: [],
  title: '',
  updatedAt: ''
}

export const emptySong: Song = {
  id: '',
  title: '',
  description: '',
  isRecorded: false,
  rehearsals: 0,
  confidence: 0,
  progress: 0,
  sections: [],
  createdAt: '',
  updatedAt: '',
  releaseDate: null,
  settings: {
    id: ''
  },
  solosCount: 0,
  riffsCount: 0
}

export const emptyPlaylist: Playlist = {
  id: '',
  title: '',
  description: '',
  songs: [],
  createdAt: '',
  updatedAt: ''
}

export const emptySongSettings: SongSettings = {
  id: ''
}

export const emptySongSection: SongSection = {
  id: '',
  name: '',
  confidence: 0,
  progress: 0,
  rehearsals: 0,
  occurrences: 0,
  partialOccurrences: 0,
  songSectionType: {
    id: '',
    name: ''
  }
}

export const emptyOrder: Order = {
  label: '',
  property: ''
}

export const defaultArtistFiltersMetadata: ArtistFiltersMetadata = {
  minBandMembersCount: 0,
  maxBandMembersCount: 5,

  minAlbumsCount: 0,
  maxAlbumsCount: 5,

  minSongsCount: 0,
  maxSongsCount: 12,

  minRehearsals: 0,
  maxRehearsals: 55,

  minConfidence: 0,
  maxConfidence: 75,

  minProgress: 0,
  maxProgress: 100
}

export const defaultAlbumFiltersMetadata: AlbumFiltersMetadata = {
  artistIds: [],

  minSongsCount: 0,
  maxSongsCount: 12,

  minRehearsals: 0,
  maxRehearsals: 55,

  minConfidence: 0,
  maxConfidence: 75,

  minProgress: 0,
  maxProgress: 100
}

export const defaultPlaylistFiltersMetadata: PlaylistFiltersMetadata = {
  minSongsCount: 0,
  maxSongsCount: 12
}
