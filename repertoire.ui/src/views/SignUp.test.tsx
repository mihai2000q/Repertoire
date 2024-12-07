import { reduxRouterRender } from '../test-utils.tsx'
import SignUp from './SignUp.tsx'
import { screen } from '@testing-library/react'

describe('Sign Up', () => {
  it('should render', () => {
    reduxRouterRender(<SignUp />)

    expect(screen.getByRole('heading', { name: /create account/i })).toBeInTheDocument()
  })
})
