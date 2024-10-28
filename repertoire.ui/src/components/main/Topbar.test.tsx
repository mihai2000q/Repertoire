import { reduxRender } from '../../test-utils.tsx'
import Topbar from './Topbar.tsx'
import { AppShell } from '@mantine/core'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import User from '../../types/models/User.ts'
import {userEvent} from "@testing-library/user-event";

describe('Topbar', () => {
  const render = () =>
    reduxRender(
      <AppShell>
        <Topbar />
      </AppShell>
    )

  const user: User = {
    id: '1',
    email: 'Gigi@yahoo.com',
    name: 'Gigi'
  }

  const handlers = [
    http.get('/users/current', async () => {
      return HttpResponse.json(user)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display search bar and user avatar', async () => {
    const [{ container }] = render()

    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument()
    expect(container.querySelector('.mantine-Loader-root')).toBeInTheDocument()
    expect(await screen.findByTestId('user-button')).toBeInTheDocument()
  })

  it('should display menu when clicking on the user button', async () => {
    // Arrange
    const userEventDispatcher = userEvent.setup()

    // Act
    render()

    const userButton = await screen.findByTestId('user-button')
    await userEventDispatcher.click(userButton)

    // Assert
    expect(screen.getByText(user.email)).toBeInTheDocument()
    expect(screen.getByText(user.name)).toBeInTheDocument()
    expect(screen.getByText('Settings')).toBeInTheDocument()
    expect(screen.getByText('Account')).toBeInTheDocument()
    expect(screen.getByText(/sign out/i)).toBeInTheDocument()
  })
})
