import { render, renderHook, RenderHookResult, RenderResult } from '@testing-library/react'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { ReactNode } from 'react'
import { Provider } from 'react-redux'
import { configureStore, EnhancedStore } from '@reduxjs/toolkit'
import { reducer, RootState } from './state/store'
import { api } from './state/api'
import { BrowserRouter } from 'react-router-dom'

export function mantineRender(ui: ReactNode) {
  return render(ui, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <MantineProvider theme={theme}>{children}</MantineProvider>
    )
  })
}

export function routerRender(ui: ReactNode) {
  return render(ui, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <BrowserRouter>
        <MantineProvider theme={theme}>{children}</MantineProvider>
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
          <MantineProvider theme={theme}>{children}</MantineProvider>
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
          <MantineProvider theme={theme}>
            {children}
          </MantineProvider>
          </BrowserRouter>
        </Provider>
      )
    }),
    store
  ]
}

// Hooks

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
