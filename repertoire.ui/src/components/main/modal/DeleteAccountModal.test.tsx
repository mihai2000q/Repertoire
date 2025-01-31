import { setupServer } from 'msw/node'
import { emptyUser, reduxRouterRender } from '../../../test-utils.tsx'
import DeleteAccountModal from './DeleteAccountModal.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { SignInRequest } from '../../../types/requests/AuthRequests.ts'
import {RootState} from "../../../state/store.ts";

describe('Delete Account Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const user = {
    ...emptyUser,
    email: 'some@email.com'
  }

  it('should render first step', () => {
    reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={() => {}}
        onCloseSettingsModal={() => {}}
        user={user}
      />
    )

    expect(screen.getByRole('dialog', { name: /delete account/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete account/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /keep account/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /continue/i })).toBeInTheDocument()
  })

  it('should render second step', async () => {
    const userEventDispatcher = userEvent.setup()

    reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={() => {}}
        onCloseSettingsModal={() => {}}
        user={user}
      />
    )

    await userEventDispatcher.click(screen.getByRole('button', { name: /continue/i }))

    expect(screen.getByRole('dialog', { name: /delete account/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete account/i })).toBeInTheDocument()
    expect(await screen.findByRole('textbox', { name: /email/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /email/i })).toHaveValue(user.email)
    expect(screen.getByRole('textbox', { name: /email/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /password/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /password/i })).toHaveValue('')
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
  })

  it('should close modals when clicking on keep account', async () => {
    const userEventDispatcher = userEvent.setup()

    const onClose = vitest.fn()

    reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={onClose}
        onCloseSettingsModal={() => {}}
        user={user}
      />
    )

    await userEventDispatcher.click(screen.getByRole('button', { name: /keep account/i }))
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should close modals when clicking on cancel', async () => {
    const userEventDispatcher = userEvent.setup()

    const onClose = vitest.fn()

    reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={onClose}
        onCloseSettingsModal={() => {}}
        user={user}
      />
    )

    await userEventDispatcher.click(screen.getByRole('button', { name: /continue/i }))
    await userEventDispatcher.click(screen.getByRole('button', { name: /cancel/i }))
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should validate password', async () => {
    const userEventDispatcher = userEvent.setup()

    reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={() => {}}
        onCloseSettingsModal={() => {}}
        user={user}
      />
    )

    await userEventDispatcher.click(screen.getByRole('button', { name: /continue/i }))
    await userEventDispatcher.click(screen.getByRole('button', { name: /delete/i }))
    expect(screen.getByText(/password cannot be blank/i)).toBeInTheDocument()
  })

  it('should send sign in request when trying to delete, to authenticate and display error if it fails', async () => {
    const userEventDispatcher = userEvent.setup()

    const password = 'somePassword'
    const error = 'Invalid credentials'

    let capturedRequest: SignInRequest
    server.use(
      http.put('/auth/sign-in', async (req) => {
        capturedRequest = (await req.request.json()) as SignInRequest
        return HttpResponse.json({ error }, { status: 401 })
      })
    )

    reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={() => {}}
        onCloseSettingsModal={() => {}}
        user={user}
      />
    )

    await userEventDispatcher.click(screen.getByRole('button', { name: /continue/i }))
    await userEventDispatcher.type(screen.getByRole('textbox', { name: /password/i }), password)
    await userEventDispatcher.click(screen.getByRole('button', { name: /delete/i }))

    expect(capturedRequest).toStrictEqual({
      email: user.email,
      password: password
    })

    expect(screen.getByText(error)).toBeInTheDocument()
  })

  it('should send sign in request and delete account requests', async () => {
    const userEventDispatcher = userEvent.setup()

    const password = 'somePassword'

    const onClose = vitest.fn()
    const onCloseSettingsModal = vitest.fn()

    let capturedRequest: SignInRequest
    server.use(
      http.put('/auth/sign-in', async (req) => {
        capturedRequest = (await req.request.json()) as SignInRequest
        return HttpResponse.json()
      }),
      http.delete('/users', async () => {
        return HttpResponse.json()
      })
    )

    const [_, store] = reduxRouterRender(
      <DeleteAccountModal
        opened={true}
        onClose={onClose}
        onCloseSettingsModal={onCloseSettingsModal}
        user={user}
      />
    )

    await userEventDispatcher.click(screen.getByRole('button', { name: /continue/i }))
    await userEventDispatcher.type(screen.getByRole('textbox', { name: /password/i }), password)
    await userEventDispatcher.click(screen.getByRole('button', { name: /delete/i }))

    expect(capturedRequest).toStrictEqual({
      email: user.email,
      password: password
    })

    expect(onClose).toHaveBeenCalledOnce()
    expect(onCloseSettingsModal).toHaveBeenCalledOnce()

    expect((store.getState() as RootState).auth.token).toBeNull()

    expect(window.location.pathname).toBe('/sign-in')
  })
})
