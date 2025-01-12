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
    expect(screen.getByRole('link', { name: /sign in/i })).toBeVisible()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeVisible()
    expect(screen.getByRole('textbox', { name: /email/i })).toBeVisible()
    expect(screen.getByLabelText(/password/i)).toBeVisible()
    expect(screen.getByRole('button', { name: /sign up/i })).toBeVisible()
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

    const passwordInput = screen.getByLabelText(/password/i)
    act(() => passwordInput.focus())
    act(() => passwordInput.blur())

    expect(screen.getByText(nameError)).toBeVisible()
    expect(screen.getByText(emailError)).toBeVisible()
    expect(screen.getByText(passwordError)).toBeVisible()
  })

  it.each([['  ', 'Name cannot be blank']])('should display name errors', async (name, error) => {
    const user = userEvent.setup()

    reduxRouterRender(<SignUp />)

    const nameInput = screen.getByRole('textbox', { name: /name/i })
    await user.type(nameInput, name)
    act(() => nameInput.blur())

    expect(screen.getByText(error)).toBeVisible()
  })

  it.each([['  '], ['email'], ['email@yahoo'], ['email.com']])(
    'should display email errors',
    async (email) => {
      const error = 'Email is invalid'
      const user = userEvent.setup()

      // Act
      reduxRouterRender(<SignUp />)

      const emailInput = screen.getByRole('textbox', { name: /email/i })
      await user.type(emailInput, email)
      act(() => emailInput.blur())

      // Assert
      expect(screen.getByText(error)).toBeVisible()
    }
  )

  it.each([
    ['1234567', 'Password must have at least 8 characters'],
    ['This is long', 'Password must have at least 1 digit'],
    ['THIS PASSWORD IS MISSING 1 LOWER CHARACTER', 'Password must have at least 1 lower character'],
    ['this password is missing 1 upper character', 'Password must have at least 1 upper character']
  ])('should display password errors', async (password, error) => {
    const user = userEvent.setup()

    reduxRouterRender(<SignUp />)

    const passwordInput = screen.getByLabelText(/password/i)
    await user.type(passwordInput, password)
    act(() => passwordInput.blur())

    expect(screen.getByText(error)).toBeVisible()
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

    screen.getAllByText(error).forEach((e) => expect(e).toBeVisible())
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

    const passwordInput = screen.getByLabelText(/password/i)
    await user.type(passwordInput, password)

    const signUpButton = screen.getByRole('button', { name: /sign up/i })
    await user.click(signUpButton)

    return store
  }
})
