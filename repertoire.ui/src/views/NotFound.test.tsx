import { mantineRender } from '../test-utils.tsx'
import NotFound from './NotFound.tsx'
import { screen } from '@testing-library/react'

describe('Not Found', () => {
  it('should render', () => {
    mantineRender(<NotFound />)

    expect(screen.getByRole('heading', { name: /not found/i })).toBeInTheDocument()
  })
})
