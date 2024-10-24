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

    expect(screen.getByText('Home')).toBeInTheDocument()
    expect(screen.getByText('Artists')).toBeInTheDocument()
    expect(screen.getByText('Albums')).toBeInTheDocument()
    expect(screen.getByText('Songs')).toBeInTheDocument()
    expect(screen.getByText('Playlists')).toBeInTheDocument()
  })

  it('should navigate to Navigation Links', async () => {
    //Arrange
    const user = userEvent.setup()

    // Act
    render()

    // Assert
    await navigateAndAssert(user, /home/i)
    await navigateAndAssert(user, /artists/i)
    await navigateAndAssert(user, /albums/i)
    await navigateAndAssert(user, /songs/i)
    await navigateAndAssert(user, /playlists/i)
  })

  async function navigateAndAssert(user: UserEvent, link: RegExp) {
    const navLink = screen.getByText(link)
    await user.click(navLink)
    expect(window.location.pathname).toMatch(link)
  }
})
