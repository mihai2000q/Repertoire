import { mantineRender } from '../test-utils.tsx'
import Home from './Home.tsx'
import { screen } from '@testing-library/react'

describe('Home', () => {
  it('should render', () => {
    mantineRender(<Home />)

    expect(screen.getByRole('heading', { name: /home/i })).toBeInTheDocument()
  })
})
