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
import Album from "./types/models/Album.ts";
import Song from "./types/models/Song.ts";
import Artist from "./types/models/Artist.ts";
import Order from "./types/Order.ts";

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

export function mantineRenderHook(hook: (props: unknown) => unknown) {
  return renderHook(hook, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <MantineProviderComponent>{children}</MantineProviderComponent>
    )
  })
}

export function reduxRenderHook(
  hook: (props: unknown) => unknown,
  preloadedState?: Partial<RootState>
): [RenderHookResult<unknown, unknown>, EnhancedStore] {
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

export function reduxRouterRenderHook(
  hook: (props: unknown) => unknown,
  preloadedState?: Partial<RootState>
): [RenderHookResult<unknown, unknown>, EnhancedStore] {
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

export const emptyArtist: Artist = {
  id: '',
  name: '',
  createdAt: '',
  updatedAt: '',
  albums: [],
  songs: []
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
  updatedAt: ''
}

export const emptyOrder: Order = {
  label: '',
  value: ''
}
