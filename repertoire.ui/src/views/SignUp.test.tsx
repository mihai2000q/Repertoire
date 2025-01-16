import { reduxRouterRender } from '../test-utils.tsx'
import SignUp from './SignUp.tsx'
import { act, screen } from '@testing-library/react'
import { expect } from 'vitest'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { SignUpRequest } from '../types/requests/AuthRequests.ts'
import TokenResponse from '../types/responses/TokenResponse.ts'
import { RootState } from '../state/store.ts'
import { setupServer } from 'msw/node'

describe('Sign Up', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRouterRender(<SignUp />)

    expect(screen.getByRole('heading', { name: /create account/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /sign in/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /email/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /password/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign up/i })).toBeInTheDocument()
  })

  it('should display validation errors when the fields are empty', async () => {
    const nameError = 'Name cannot be blank'
    const emailError = 'Email is invalid'
    const passwordError = 'Password cannot be blank'

    reduxRouterRender(<SignUp />)

    const nameInput = screen.getByRole('textbox', { name: /name/i })
    act(() => nameInput.focus())
    act(() => nameInput.blur())

    const emailInput = screen.getByRole('textbox', { name: /email/i })
    act(() => emailInput.focus())
    act(() => emailInput.blur())

    const passwordInput = screen.getByRole('textbox', { name: /password/i })
    act(() => passwordInput.focus())
    act(() => passwordInput.blur())

    expect(screen.getByText(nameError)).toBeInTheDocument()
    expect(screen.getByText(emailError)).toBeInTheDocument()
    expect(screen.getByText(passwordError)).toBeInTheDocument()
  })

  it.each([['  ', /name cannot be blank/i]])('should display name errors', async (name, error) => {
    const user = userEvent.setup()

    reduxRouterRender(<SignUp />)

    const nameInput = screen.getByRole('textbox', { name: /name/i })
    await user.type(nameInput, name)
    act(() => nameInput.blur())
    expect(nameInput).toBeInvalid()
    expect(screen.getByText(error)).toBeInTheDocument()
  })

  it.each([['  '], ['email'], ['email@yahoo'], ['email.com']])(
    'should display email errors',
    async (email) => {
      const error = /email is invalid/i
      const user = userEvent.setup()

      reduxRouterRender(<SignUp />)

      const emailInput = screen.getByRole('textbox', { name: /email/i })
      await user.type(emailInput, email)
      act(() => emailInput.blur())
      expect(emailInput).toBeInvalid()
      expect(screen.getByText(error)).toBeInTheDocument()
    }
  )

  it.each([
    ['1234567', /password must have at least 8 characters/i],
    ['This is long', /password must have at least 1 digit/i],
    [
      'THIS PASSWORD IS MISSING 1 LOWER CHARACTER',
      /password must have at least 1 lower character/i
    ],
    ['this password is missing 1 upper character', /Password must have at least 1 upper character/i]
  ])('should display password errors', async (password, error) => {
    const user = userEvent.setup()

    reduxRouterRender(<SignUp />)

    const passwordInput = screen.getByRole('textbox', { name: /password/i })
    await user.type(passwordInput, password)
    act(() => passwordInput.blur())
    expect(screen.getByText(error)).toBeInTheDocument()
    expect(passwordInput).toHaveAttribute('data-invalid', 'true')
  })

  it('should send sign up request and display sign up error', async () => {
    const name = 'Someone Else'
    const email = 'someone@else.com'
    const password = 'ThisIsAGoodPassword123'
    const error = 'Email already in use'

    server.use(
      http.post('/auth/sign-up', async () => HttpResponse.json({ error }, { status: 401 }))
    )

    await sendSignUpRequest(name, email, password)

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()
    expect(screen.getByRole('textbox', { name: /email/i })).toBeInvalid()
    expect(screen.getByRole('textbox', { name: /password/i })).toHaveAttribute(
      'data-invalid',
      'true'
    )
    expect(screen.getAllByText(error)).toHaveLength(3)
  })

  it('should send sign up request and save token', async () => {
    const name = 'Someone Else'
    const email = 'someone@else.com'
    const password = 'ThisIsAGoodPassword123'

    let capturedSignUpRequest: SignUpRequest | undefined

    const expectedToken = 'token'

    server.use(
      http.post('/auth/sign-up', async (req) => {
        capturedSignUpRequest = (await req.request.json()) as SignUpRequest
        const response: TokenResponse = { token: expectedToken }
        return HttpResponse.json(response)
      })
    )

    const store = await sendSignUpRequest(name, email, password)

    expect(capturedSignUpRequest).toStrictEqual({ name, email, password })
    expect((store.getState() as RootState).auth.token).toBe(expectedToken)
    expect(window.location.pathname).toBe('/home')
  })

  async function sendSignUpRequest(name: string, email: string, password: string) {
    const user = userEvent.setup()

    const [_, store] = reduxRouterRender(<SignUp />)

    const nameInput = screen.getByRole('textbox', { name: /name/i })
    await user.type(nameInput, name)

    const emailInput = screen.getByRole('textbox', { name: /email/i })
    await user.type(emailInput, email)

    const passwordInput = screen.getByRole('textbox', { name: /password/i })
    await user.type(passwordInput, password)

    const signUpButton = screen.getByRole('button', { name: /sign up/i })
    await user.click(signUpButton)

    return store
  }
})
