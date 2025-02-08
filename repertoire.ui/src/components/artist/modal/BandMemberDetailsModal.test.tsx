import { mantineRender } from '../../../test-utils.tsx'
import BandMemberDetailsModal from './BandMemberDetailsModal.tsx'
import { BandMember } from '../../../types/models/Artist.ts'
import { screen } from '@testing-library/react'

describe('Band Member Details Modal', () => {
  const bandMember: BandMember = {
    id: '1',
    name: 'Member',
    color: '#123123',
    roles: [
      { id: '2', name: 'Some Role' },
      { id: '3', name: 'Some Other Role' }
    ]
  }

  it('should render', () => {
    mantineRender(
      <BandMemberDetailsModal opened={true} onClose={() => {}} bandMember={bandMember} />
    )

    expect(screen.getByRole('dialog')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: bandMember.name })).toBeInTheDocument()
    bandMember.roles.forEach((role) => expect(screen.getByText(role.name)).toBeInTheDocument())
  })

  it('should render with image', () => {
    const image = 'something.png'

    mantineRender(
      <BandMemberDetailsModal
        opened={true}
        onClose={() => {}}
        bandMember={{ ...bandMember, imageUrl: image }}
      />
    )

    expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: bandMember.name })).toHaveAttribute('src', image)
  })
})
