import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import BandMembersGroup from './BandMembersGroup.tsx'
import { BandMember } from '../../../../../types/models/Artist.ts'
import { mantineRender } from '../../../../../test-utils.tsx'

const bandMembers: BandMember[] = [
  {
    id: '1',
    name: 'John Doe',
    color: 'blue',
    imageUrl: '',
    roles: []
  },
  {
    id: '2',
    name: 'Jane Smith',
    color: 'red',
    imageUrl: 'something.png',
    roles: []
  }
]

describe('Band Members Group', () => {
  it('should render an avatar per band member', () => {
    mantineRender(<BandMembersGroup bandMembers={bandMembers} />)

    bandMembers.forEach((bandMember) => {
      if (bandMember.imageUrl) {
        expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
      } else {
        expect(screen.getByLabelText(`default-icon-${bandMember.name}`)).toBeInTheDocument()
      }
    })
  })

  it('should show all band member names in the tooltip on hover', async () => {
    const user = userEvent.setup()

    mantineRender(<BandMembersGroup bandMembers={bandMembers} />)

    await user.hover(screen.getByRole('img', { name: bandMembers[1].name }))

    const tooltip = await screen.findByRole('tooltip')

    bandMembers.forEach((bandMember) => {
      expect(tooltip).toHaveTextContent(bandMember.name)
    })
  })
})
