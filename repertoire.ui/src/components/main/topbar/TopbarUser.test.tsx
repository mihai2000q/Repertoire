import User from '../../../types/models/User.ts'
import { emptyUser, reduxRouterRender } from '../../../test-utils.tsx'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { RootState } from '../../../state/store.ts'
import TopbarUser from './TopbarUser.tsx'

describe('Topbar User', () => {
  const render = () => reduxRouterRender(<TopbarUser />, { auth: { token: 'some token' } })

  const user: User = {
    ...emptyUser,
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

  it('should display menu when clicking on the user button', async () => {
    const userEventDispatcher = userEvent.setup()

    const [_, store] = render()

    const userButton = await screen.findByRole('button', { name: 'user' })
    await userEventDispatcher.click(userButton)

    expect(screen.getByText(user.email)).toBeInTheDocument()
    expect(screen.getByText(user.name)).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /settings/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /account/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /sign out/i })).toBeInTheDocument()
    expect((store.getState() as RootState).global.userId).toBe(user.id)
  })

  describe('on menu', () => {
    it('should display account modal when clicking on account', async () => {
      const userEventDispatcher = userEvent.setup()

      render()

      await userEventDispatcher.click(await screen.findByRole('button', { name: 'user' }))
      await userEventDispatcher.click(screen.getByRole('menuitem', { name: /account/i }))

      expect(await screen.findByRole('dialog', { name: /account/i })).toBeInTheDocument()
    })

    it('should display settings modal when clicking on settings', async () => {
      const userEventDispatcher = userEvent.setup()

      render()

      await userEventDispatcher.click(await screen.findByRole('button', { name: 'user' }))
      await userEventDispatcher.click(screen.getByRole('menuitem', { name: /settings/i }))

      expect(await screen.findByRole('dialog', { name: /settings/i })).toBeInTheDocument()
    })

    it('should sign out when clicking on sign out', async () => {
      const userEventDispatcher = userEvent.setup()

      const [_, store] = render()

      await userEventDispatcher.click(await screen.findByRole('button', { name: 'user' }))
      await userEventDispatcher.click(screen.getByRole('menuitem', { name: /sign out/i }))

      expect((store.getState() as RootState).auth.token).toBeNull()
    })
  })
})
