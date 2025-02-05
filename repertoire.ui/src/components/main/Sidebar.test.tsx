import { routerRender } from '../../test-utils.tsx'
import Sidebar from './Sidebar.tsx'
import { screen } from '@testing-library/react'
import { AppShell } from '@mantine/core'
import { UserEvent, userEvent } from '@testing-library/user-event'

describe('Sidebar', () => {
  const render = () => {
    routerRender(
      <AppShell>
        <Sidebar />
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
})
