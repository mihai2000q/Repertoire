import { reduxRouterRender } from '../test-utils.tsx'
import SignIn from './SignIn.tsx'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { SignInRequest } from '../types/requests/AuthRequests.ts'
import TokenResponse from '../types/responses/TokenResponse.ts'
import { expect } from 'vitest'
import { RootState } from '../state/store.ts'

describe('Sign In', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    const [_, store] = reduxRouterRender(<SignIn />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/sign in/i)
    expect(screen.getByRole('heading', { name: /welcome/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /create account/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /email/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /password/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /forgot password/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('should display validation errors when the fields are empty', async () => {
    const emailError = 'Email is invalid'
    const passwordError = 'Password cannot be blank'

    reduxRouterRender(<SignIn />)

    const emailInput = screen.getByRole('textbox', { name: /email/i })
    act(() => emailInput.focus())
    act(() => emailInput.blur())

    const passwordInput = screen.getByRole('textbox', { name: /password/i })
    act(() => passwordInput.focus())
    act(() => passwordInput.blur())

    expect(emailInput).toBeInvalid()
    expect(passwordInput).toHaveAttribute('data-invalid', 'true')
    expect(screen.getByText(emailError)).toBeInTheDocument()
    expect(screen.getByText(passwordError)).toBeInTheDocument()
  })

  it.each([['  '], ['email'], ['email@yahoo'], ['email.com']])(
    'should display email errors',
    async (email) => {
      const error = /email is invalid/i
      const user = userEvent.setup()

      reduxRouterRender(<SignIn />)

      const emailInput = screen.getByRole('textbox', { name: /email/i })
      await user.type(emailInput, email)
      act(() => emailInput.blur())
      expect(emailInput).toBeInvalid()
      expect(screen.getByText(error)).toBeInTheDocument()
    }
  )

  it('should send sign in request and display sign in error', async () => {
    const email = 'someone@else.com'
    const password = 'ThisIsAGoodPassword123'
    const error = 'Invalid credentials'

    server.use(http.put('/auth/sign-in', async () => HttpResponse.json({ error }, { status: 401 })))

    await sendSignInRequest(email, password)

    expect(screen.getByRole('textbox', { name: /email/i })).toBeInvalid()
    expect(screen.getByRole('textbox', { name: /password/i })).toHaveAttribute(
      'data-invalid',
      'true'
    )
    expect(screen.getAllByText(error)).toHaveLength(2)
  })

  it('should send sign in request and save token', async () => {
    const email = 'someone@else.com'
    const password = 'ThisIsAGoodPassword123'

    let capturedSignInRequest: SignInRequest | undefined

    const expectedToken = 'token'

    server.use(
      http.put('/auth/sign-in', async (req) => {
        capturedSignInRequest = (await req.request.json()) as SignInRequest
        const response: TokenResponse = { token: expectedToken }
        return HttpResponse.json(response)
      })
    )

    const store = await sendSignInRequest(email, password)

    expect(capturedSignInRequest).toStrictEqual({ email, password })
    expect((store.getState() as RootState).auth.token).toBe(expectedToken)
    expect(window.location.pathname).toBe('/home')
  })

  async function sendSignInRequest(email: string, password: string) {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<SignIn />)

    const emailInput = screen.getByRole('textbox', { name: /email/i })
    await user.type(emailInput, email)

    const passwordInput = screen.getByRole('textbox', { name: /password/i })
    await user.type(passwordInput, password)

    const signInButton = screen.getByRole('button', { name: /sign in/i })
    await user.click(signInButton)

    return store
  }
})
