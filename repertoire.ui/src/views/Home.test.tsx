import { reduxRender } from '../test-utils.tsx'
import Home from './Home.tsx'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import { RootState } from '../state/store.ts'

describe('Home', () => {
  it('should render', () => {
    const [_, store] = reduxRender(<Home />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/home/i)
    expect(screen.getByRole('heading', { name: /home/i })).toBeInTheDocument()
  })
})
