import { reduxRouterRender } from '../../test-utils.tsx'
import Topbar from './Topbar.tsx'
import { AppShell } from '@mantine/core'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import User from '../../types/models/User.ts'
import { userEvent } from '@testing-library/user-event'
import {RootState} from "../../state/store.ts";

describe('Topbar', () => {
  const render = (token: string | null = 'some token') =>
    reduxRouterRender(
      <AppShell>
        <Topbar />
      </AppShell>,
      { auth: { token } }
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

  it.each(['some token', undefined])(
    'should render and display search bar and user avatar',
    async (token) => {
      render(token)

      expect(screen.getByPlaceholderText('Search')).toBeInTheDocument()
      expect(await screen.findByRole('button', { name: 'user' })).toBeInTheDocument()
    }
  )

  it('should display menu when clicking on the user button', async () => {
    // Arrange
    const userEventDispatcher = userEvent.setup()

    // Act
    render()

    const userButton = await screen.findByRole('button', { name: 'user' })
    await userEventDispatcher.click(userButton)

    // Assert
    expect(screen.getByText(user.email)).toBeInTheDocument()
    expect(screen.getByText(user.name)).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /settings/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /account/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /sign out/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display account modal when clicking on account', async () => {
      // Arrange
      const userEventDispatcher = userEvent.setup()

      // Act
      render()

      await userEventDispatcher.click(await screen.findByRole('button', { name: 'user' }))
      await userEventDispatcher.click(screen.getByRole('menuitem', { name: /account/i }))

      // Assert
      expect(screen.getByRole('heading', { name: /account/i })).toBeInTheDocument()
    })

    it('should sign out when clicking on sign out', async () => {
      // Arrange
      const userEventDispatcher = userEvent.setup()

      // Act
      const [_, store] = render()

      await userEventDispatcher.click(await screen.findByRole('button', { name: 'user' }))
      await userEventDispatcher.click(screen.getByRole('menuitem', { name: /sign out/i }))

      // Assert
      expect((store.getState() as RootState).auth.token).toBeNull()
    })
  })
})
