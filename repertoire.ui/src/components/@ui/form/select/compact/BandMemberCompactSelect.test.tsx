import { reduxRender } from '../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { BandMember } from '../../../../../types/models/Artist.ts'
import BandMemberCompactSelect from './BandMemberCompactSelect.tsx'

describe('Band Member Compact Select', () => {
  const bandMembers: BandMember[] = [
    {
      id: '1',
      name: 'Chester',
      roles: [{ id: '1', name: 'Vocals' }]
    },
    {
      id: '2',
      name: 'Michael',
      roles: [
        { id: '2', name: 'Vocals' },
        { id: '3', name: 'Guitarist' }
      ]
    },
    {
      id: '3',
      name: 'Luther',
      roles: [{ id: '1', name: 'Bassist' }]
    }
  ]

  it('should render and change band members', async () => {
    const user = userEvent.setup()

    const newBandMember = bandMembers[0]

    const setBandMember = vitest.fn()

    const [{ rerender }] = reduxRender(
      <BandMemberCompactSelect
        bandMember={null}
        setBandMember={setBandMember}
        bandMembers={bandMembers}
      />
    )

    const selectButton = screen.getByRole('button', { name: 'select-band-member' })
    expect(selectButton).not.toBeDisabled()

    await user.hover(selectButton)
    expect(await screen.findByRole('tooltip', { name: /choose a band member/i })).toBeInTheDocument()

    await user.click(selectButton)

    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue('')

    for (const member of bandMembers) {
      expect(await screen.findByRole('option', { name: member.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newBandMember.name })
    await user.click(selectedOption)

    expect(setBandMember).toHaveBeenCalledOnce()
    expect(setBandMember).toHaveBeenCalledWith(newBandMember)

    rerender(
      <BandMemberCompactSelect
        bandMember={newBandMember}
        setBandMember={setBandMember}
        bandMembers={bandMembers}
      />
    )

    expect(screen.queryByRole('button', { name: 'select-band-member' })).not.toBeInTheDocument()

    const memberButton = screen.getByRole('button', { name: newBandMember.name })
    expect(memberButton).toBeInTheDocument()

    await user.hover(memberButton)
    expect(await screen.findByRole('tooltip', { name: new RegExp(newBandMember.name, 'i') })).toBeInTheDocument()

    await user.click(memberButton)
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue(newBandMember.name)

    // reset the value from outside component
    rerender(
      <BandMemberCompactSelect
        bandMember={null}
        setBandMember={setBandMember}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'select-band-member' }))
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue('')
  })

  it('should filter by name', async () => {
    const user = userEvent.setup()

    const searchValue = 't'

    reduxRender(
      <BandMemberCompactSelect
        bandMember={null}
        setBandMember={() => {}}
        bandMembers={bandMembers}
      />
    )

    await user.click(screen.getByRole('button', { name: 'select-band-member' }))
    await user.type(screen.getByRole('textbox', { name: /search/i }), searchValue)

    const filteredMembers = bandMembers.filter((b) => b.name.includes(searchValue))
    expect(await screen.findAllByRole('option')).toHaveLength(filteredMembers.length)
    for (const member of filteredMembers) {
      expect(screen.getByRole('option', { name: member.name })).toBeInTheDocument()
    }
  })

  it('should be disabled when the band members are undefined', async () => {
    const user = userEvent.setup()

    reduxRender(
      <BandMemberCompactSelect bandMember={null} setBandMember={() => {}} bandMembers={undefined} />
    )

    const button = screen.getByRole('button', { name: 'select-band-member' })
    expect(button).toBeDisabled()
    await user.hover(button)
    expect(await screen.findByRole('tooltip', { name: /not a band/i })).toBeInTheDocument()
  })
})
