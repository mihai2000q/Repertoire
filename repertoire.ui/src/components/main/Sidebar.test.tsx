import { routerRender } from '../../test-utils.tsx'
import Sidebar from './Sidebar.tsx'
import { screen } from '@testing-library/react'
import { AppShell } from '@mantine/core'
import { UserEvent, userEvent } from '@testing-library/user-event'

describe('Sidebar', () => {
  const render = (toggleSidebar: () => void = () => {}) => {
    routerRender(
      <AppShell
        navbar={{ width: 200, breakpoint: 'sm', collapsed: { mobile: false, desktop: false } }}
      >
        <Sidebar toggleSidebar={toggleSidebar} />
      </AppShell>
    )
  }

  it('should render and display Navigation Links', ({ expect }) => {
    render()

    expect(screen.getByRole('link', { name: /home/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /artists/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /albums/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /playlists/i })).toBeInTheDocument()
  })

  it('should navigate to Navigation Links', async () => {
    //Arrange
    const user = userEvent.setup()

    render()

    await navigateAndAssert(user, /home/i)
    await navigateAndAssert(user, /artists/i)
    await navigateAndAssert(user, /albums/i)
    await navigateAndAssert(user, /songs/i)
    await navigateAndAssert(user, /playlists/i)
  })

  async function navigateAndAssert(user: UserEvent, link: RegExp) {
    const navLink = screen.getByRole('link', { name: link })
    await user.click(navLink)
    expect(window.location.pathname).toMatch(link)
    expect(navLink).toHaveAttribute('data-active', 'true')
  }

  it.skip('should display navigation menu button, when screen is small', async () => {
    const userEventDispatcher = userEvent.setup()

    const originalInnerWidth = window.innerWidth
    vi.spyOn(window, 'innerWidth', 'get').mockImplementation(() => 500)

    const toggleSidebar = vitest.fn()

    render(toggleSidebar)

    const button = screen.getByRole('button', { name: 'toggle-sidebar' })
    expect(button).toBeInTheDocument()
    await userEventDispatcher.click(button)
    expect(toggleSidebar).toHaveBeenCalledOnce()

    // restore
    vi.spyOn(window, 'innerWidth', 'get').mockImplementation(() => originalInnerWidth)
  })
})
