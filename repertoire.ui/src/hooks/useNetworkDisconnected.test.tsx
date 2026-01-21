import { withToastify } from '../test-utils.tsx'
import useNetworkDisconnected from './useNetworkDisconnected.tsx'
import { act, renderHook, screen } from '@testing-library/react'
import { ReactNode } from 'react'
import { theme } from '../theme/theme.ts'
import { emotionTransform, MantineEmotionProvider } from '@mantine/emotion'
import { MantineProvider } from '@mantine/core'

describe('use Network Disconnected', () => {
  it('should render and display toasts', async () => {
    // mock the navigator
    const originalNavigator = { ...navigator }

    Object.defineProperty(window, 'navigator', {
      writable: true,
      value: {
        ...originalNavigator,
        online: true
      }
    })

    renderHook(() => useNetworkDisconnected(), {
      wrapper: ({ children }: { children: ReactNode }) => (
        <MantineProvider theme={theme} stylesTransform={emotionTransform}>
          <MantineEmotionProvider>{withToastify(children)}</MantineEmotionProvider>
        </MantineProvider>
      )
    })

    // go offline
    act(() => {
      Object.defineProperty(window.navigator, 'online', {
        writable: true,
        value: false
      })
      window.dispatchEvent(new Event('offline'))
    })

    expect(await screen.findByText(/please check your internet connection/i)).toBeInTheDocument()

    // go back online
    act(() => {
      Object.defineProperty(window.navigator, 'online', {
        writable: true,
        value: true
      })
      window.dispatchEvent(new Event('online'))
    })

    expect(await screen.findByText(/internet is back/i)).toBeInTheDocument()

    // restore
    Object.defineProperty(window, 'navigator', {
      writable: true,
      value: originalNavigator
    })

    vi.restoreAllMocks()
  })
})
