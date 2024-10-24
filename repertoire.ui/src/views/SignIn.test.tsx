import { reduxRouterRender } from '../test-utils.tsx'
import SignIn from './SignIn.tsx'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import SignInRequest from '../types/requests/AuthRequests.ts'
import TokenResponse from '../types/responses/TokenResponse.ts'
import {EnhancedStore} from "@reduxjs/toolkit";
import {expect} from "vitest";
import {RootState} from "../state/store.ts";

describe('Sign In', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display specific elements', () => {
    const [{ container }] = reduxRouterRender(<SignIn />)

    expect(screen.getByRole('heading', { name: /welcome/i })).toBeVisible()
    expect(screen.getByRole('textbox', { name: /email/i })).toBeVisible()
    expect(container.querySelector('input[type=password]')).toBeVisible()
    expect(screen.getByRole('button', { name: /create account/i })).toBeVisible()
    expect(screen.getByRole('button', { name: /forgot password/i })).toBeVisible()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeVisible()
  })

  it('should display validation errors when the fields are empty', async () => {
    // Arrange
    const emailError = 'Email is invalid'
    const passwordError = 'Password cannot be blank'

    // Act
    const [{ container }] = reduxRouterRender(<SignIn />)

    const emailInput = screen.getByRole('textbox', { name: /email/i })
    act(() => emailInput.focus())
    act(() => emailInput.blur())

    const passwordInput = container.querySelector('input[type=password]') as HTMLElement
    act(() => passwordInput.focus())
    act(() => passwordInput.blur())

    // Assert
    expect(screen.getByText(emailError)).toBeVisible()
    expect(screen.getByText(passwordError)).toBeVisible()
  })

  it.each([['  '], ['email'], ['email@yahoo'], ['email.com']])(
    'should display email errors',
    async (email) => {
      // Arrange
      const error = 'Email is invalid'
      const user = userEvent.setup()

      // Act
      reduxRouterRender(<SignIn />)

      const emailInput = screen.getByRole('textbox', { name: /email/i })
      await user.type(emailInput, email)
      act(() => emailInput.blur())

      // Assert
      expect(screen.getByText(error)).toBeVisible()
    }
  )

  it('should display password errors', async () => {
    // Arrange
    const user = userEvent.setup()
    const error = 'Password cannot be blank'

    // Act
    const [{ container }] = reduxRouterRender(<SignIn />)

    const passwordInput = container.querySelector('input[type=password]') as HTMLElement
    await user.type(passwordInput, ' ')
    act(() => passwordInput.blur())

    // Assert
    expect(screen.getByText(error)).toBeVisible()
  })

  it('should send sign in request and display sign in error', async () => {
    // Arrange
    const email = 'someone@else.com'
    const password = 'ThisIsAGoodPassword123'
    const error = 'Invalid credentials'

    server.use(
      http.put('/auth/sign-in', async () => HttpResponse.json({ error }, { status: 401 }))
    )

    // Act
    await sendSignInRequest(email, password)

    // Assert
    screen.getAllByText(error).forEach(e => expect(e).toBeVisible())
  })

  it('should send sign in request and save token', async () => {
    // Arrange
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

    // Act
    const store = await sendSignInRequest(email, password)

    // Assert
    expect(capturedSignInRequest).toStrictEqual({ email, password })
    expect((store.getState() as RootState).auth.token).toBe(expectedToken)
    expect(window.location.pathname).toBe('/home')
  })

  async function sendSignInRequest(email: string, password: string): Promise<EnhancedStore> {
    const user = userEvent.setup()

    const [{ container }, store] = reduxRouterRender(<SignIn />)

    const emailInput = screen.getByRole('textbox', { name: /email/i })
    await user.type(emailInput, email)

    const passwordInput = container.querySelector('input[type=password]') as HTMLElement
    await user.type(passwordInput, password)

    const signInButton = screen.getByRole('button', { name: /sign in/i })
    await user.click(signInButton)

    return store
  }
})
