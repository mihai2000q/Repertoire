import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { BandMember } from '../../../../types/models/Artist.ts'
import BandMemberSelect from './BandMemberSelect.tsx'

describe('Band Member Select', () => {
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

  it('should render and band members', async () => {
    const user = userEvent.setup()

    const newBandMember = bandMembers[0]

    const setBandMember = vitest.fn()

    const [{ rerender }] = reduxRender(
      <BandMemberSelect bandMember={null} setBandMember={setBandMember} bandMembers={bandMembers} />
    )

    const select = screen.getByRole('textbox', { name: /band member/i })
    expect(select).toHaveValue('')
    await user.click(select)

    for (const member of bandMembers) {
      expect(await screen.findByRole('option', { name: member.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newBandMember.name })
    await user.click(selectedOption)

    expect(setBandMember).toHaveBeenCalledOnce()
    expect(setBandMember).toHaveBeenCalledWith(newBandMember)

    rerender(
      <BandMemberSelect
        bandMember={newBandMember}
        setBandMember={setBandMember}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('textbox', { name: /band member/i })).toHaveValue(newBandMember.name)
  })
})
