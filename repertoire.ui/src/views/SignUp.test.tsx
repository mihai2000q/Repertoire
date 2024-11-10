import { mantineRender } from '../test-utils.tsx'
import SignUp from './SignUp.tsx'
import { screen } from '@testing-library/react'

describe('Sign Up', () => {
  it('should render', () => {
    mantineRender(<SignUp />)

    expect(screen.getByRole('heading', { name: /sign up/i })).toBeInTheDocument()
  })
})
