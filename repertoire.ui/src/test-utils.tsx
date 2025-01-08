import { render, renderHook, RenderHookResult, RenderResult } from '@testing-library/react'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { ReactNode } from 'react'
import { Provider } from 'react-redux'
import { configureStore, EnhancedStore } from '@reduxjs/toolkit'
import { reducer, RootState } from './state/store'
import { api } from './state/api'
import { BrowserRouter, MemoryRouter, Route, Routes } from 'react-router-dom'
import { emotionTransform, MantineEmotionProvider } from '@mantine/emotion'
import { ToastContainer } from 'react-toastify'

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
  const store = configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })

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
  const store = configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })

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
  const store = configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })

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
  const store = configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })

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
  const store = configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })

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
