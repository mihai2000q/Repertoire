import { render as testingLibraryRender } from '@testing-library/react'
import { MantineProvider } from '@mantine/core'
import { theme } from './theme/theme'
import { ReactNode } from 'react'
import { Provider } from 'react-redux'
import { configureStore } from '@reduxjs/toolkit'
import { reducer, RootState } from './state/store'
import {api} from "./state/api";

export function mantineRender(ui: ReactNode) {
  return testingLibraryRender(ui, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <MantineProvider theme={theme}>{children}</MantineProvider>
    )
  })
}

export function reduxRender(ui: ReactNode, preloadedState?: Partial<RootState>) {
  const store = configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })

  return testingLibraryRender(ui, {
    wrapper: ({ children }: { children: ReactNode }) => (
      <Provider store={store}>
        <MantineProvider theme={theme}>{children}</MantineProvider>
      </Provider>
    )
  })
}
