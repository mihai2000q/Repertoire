import { reduxRouterRender } from '../../../test-utils.tsx'
import TopbarNavigation from './TopbarNavigation.tsx'
import { act, screen } from '@testing-library/react'
import { signIn } from '../../../state/slice/authSlice.ts'
import { RootState } from '../../../state/store.ts'

describe('Topbar Navigation', () => {
  afterEach(() => (location.pathname = '/'))

  it('should render', () => {
    const [{ rerender }, store] = reduxRouterRender(<TopbarNavigation />, {
      auth: { token: '', historyOnSignIn: { index: 0, justSignedIn: false } }
    })

    const backButton = screen.getByRole('button', { name: 'back' })
    const forwardButton = screen.getByRole('button', { name: 'forward' })

    // Initial render
    expect(backButton).toBeDisabled()
    expect(forwardButton).toBeDisabled()

    // Simulate double forward navigation
    act(() => {
      history.pushState({ idx: 1 }, '', '/test-1')
      history.pushState({ idx: 2 }, '', '/test-2')
    })
    rerender(<TopbarNavigation />)

    // should enable back button on forward navigation
    expect(backButton).not.toBeDisabled()
    expect(forwardButton).toBeDisabled()

    // Simulate back navigation
    act(() => history.replaceState({ idx: 1 }, '', '/test-1'))
    rerender(<TopbarNavigation />)

    // Should enable forward navigation when back navigation
    expect(forwardButton).not.toBeDisabled()

    // Should disable forward navigation when sign-in and back button by last index
    act(() => store.dispatch(signIn('some-token')))
    expect((store.getState() as RootState).auth.historyOnSignIn).toStrictEqual({
      index: 2,
      justSignedIn: true
    })
    expect(backButton).toBeDisabled()
    expect(forwardButton).toBeDisabled()

    act(() => history.pushState({ idx: 3 }, '', '/test-3'))
    rerender(<TopbarNavigation />)
    expect((store.getState() as RootState).auth.historyOnSignIn).toStrictEqual({
      index: 2,
      justSignedIn: false
    })
    expect(backButton).not.toBeDisabled()

    // Simulate back navigation
    act(() => history.replaceState({ idx: 1 }, '', '/test-1'))
    rerender(<TopbarNavigation />)

    expect(screen.getByRole('button', { name: 'forward' })).not.toBeDisabled()
  })
})
