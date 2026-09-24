import { mantineRender } from '../../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import BandMemberAvatar from './BandMemberAvatar.tsx'
import { BandMember } from '../../../../../../types/models/Artist.ts'

const bandMember: BandMember = {
  id: '1',
  name: 'John Doe',
  color: 'blue',
  imageUrl: 'something.png',
  roles: [{ id: 'r1', name: 'Guitarist' }]
}

describe('Band Member Avatar', () => {
  it('should render avatar when there is an image', () => {
    mantineRender(<BandMemberAvatar bandMember={bandMember} />)

    expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
  })

  it('should render default icon when there is no image', () => {
    const localBandMember = { ...bandMember, imageUrl: '' }

    mantineRender(<BandMemberAvatar bandMember={localBandMember} />)

    expect(screen.getByLabelText(`default-icon-${bandMember.name}`)).toBeInTheDocument()
  })

  it('should show details on hover', async () => {
    const user = userEvent.setup()


    mantineRender(<BandMemberAvatar bandMember={bandMember} />)

    expect(screen.queryByText(bandMember.roles[0].name)).not.toBeInTheDocument()
    await user.hover(screen.getByRole('img', { name: bandMember.name }))
    expect(await screen.findByText(bandMember.roles[0].name)).toBeInTheDocument()
    expect(screen.getAllByRole('img', { name: bandMember.name })).toHaveLength(2)
  })

  it('should not append ellipsis when there are 2 or fewer roles', async () => {
    const user = userEvent.setup()
    const localBandMember: BandMember = {
      ...bandMember,
      roles: [
        { id: 'r1', name: 'Guitarist' },
        { id: 'r2', name: 'Singer' }
      ]
    }

    const firstRole = localBandMember.roles[0].name
    const secondRole = localBandMember.roles[1].name

    mantineRender(<BandMemberAvatar bandMember={localBandMember} />)
    await user.hover(screen.getByRole('img', { name: localBandMember.name }))

    expect(await screen.findByText(firstRole)).toBeInTheDocument()
    expect(screen.getByText(secondRole)).toBeInTheDocument()
    expect(screen.queryByText(new RegExp(`^${secondRole} \\.\\.\\.`))).not.toBeInTheDocument()
  })

  it('should show up to 2 roles with ellipsis when more exist', async () => {
    const user = userEvent.setup()
    const localBandMember = {
      ...bandMember,
      roles: [
        { id: 'r1', name: 'Guitarist' },
        { id: 'r2', name: 'Singer' },
        { id: 'r3', name: 'Songwriter' }
      ]
    }

    const firstRole = localBandMember.roles[0].name
    const secondRole = localBandMember.roles[1].name
    const thirdRole = localBandMember.roles[2].name

    mantineRender(<BandMemberAvatar bandMember={localBandMember} />)
    await user.hover(screen.getByRole('img', { name: localBandMember.name }))

    expect(await screen.findByText(firstRole)).toBeInTheDocument()
    expect(screen.getByText(new RegExp(`^${secondRole}`))).toHaveTextContent(`${secondRole} ...`)
    expect(screen.queryByText(thirdRole)).not.toBeInTheDocument()
  })
})
