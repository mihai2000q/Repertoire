import { mantineRender } from '../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import RehearsalsBadge from './RehearsalsBadge.tsx'

describe('Rehearsals Badge', () => {
  it('should render', () => {
    const rehearsals = 7

    mantineRender(<RehearsalsBadge rehearsals={rehearsals} />)

    expect(screen.getByText(rehearsals.toString())).toBeInTheDocument()
  })
})
