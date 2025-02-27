%lang starknet

from starkware.cairo.common.bitwise import bitwise_and
from starkware.cairo.common.cairo_builtins import (HashBuiltin,
    BitwiseBuiltin)
from starkware.cairo.common.math import (unsigned_div_rem,
    assert_not_zero)
from starkware.cairo.common.pow import pow
from contracts.utils.interfaces import IModuleController
from starkware.starknet.common.syscalls import get_caller_address

##### Module 04 #####
#
# This module is populated using NFT ownership data from other
# chains (DOPE-L1 or Hustler-Optimism). Once 'imported' and stored
# here, modules can then use that information. Likely going to
# involve L1 storage proofs (WIP).
#
# Currently has scores feed from off-chain for some items,
# which may be replaced by Module 10 scoring
# system.
#
####################

##### Encoding details #####
# Zero-based bit index for data locations.
# 0 weapon id.
# 6 weapon strength score (v1).
# 10 clothes.
# 20 vehicle id.
# 26 vehicle speed score (v1).
# 30 waistArmor id.
# 40 footArmor id.
# 46 footArmor speed score (v1).
# 50 handArmor id.
# 60 necklace id.
# 66 necklace bribe score (v1).
# 70 ring id.
# 76 ring bribe score (v1).
# 80 suffix id.
# 90 drug id (v1).
# 100 namePrefixes.
# 110 nameSuffixes.
# 120-249 (vacant).

# Test data with alternating id/score values: 113311331133113311331133
# E.g., weapon score = 3, vehicle speed score = 3, ring bribe score = 1.
# 00010000010011000011 * 6 = 12 items (indices starting 0-110).
#000100000100110000110001000001001100001100010000010011000011000100000100110000110001000001001100001100010000010011000011
const TESTDATA1 = 84622096520155505419920978765481155

##### Storage #####
# Binary encoding of ownership fields.
@storage_var
func user_data(
        user_id : felt
    ) -> (
        data : felt
    ):
end

@storage_var
func user_count(
    ) -> (
        res : felt
    ):
end


@storage_var
func controller_address() -> (address : felt):
end


# Called on deployment only.
@constructor
func constructor{
        syscall_ptr : felt*,
        pedersen_ptr : HashBuiltin*,
        range_check_ptr
    }(
        address_of_controller : felt
    ):
    # Store the address of the only fixed contract in the system.
    controller_address.write(address_of_controller)
    return ()
end

############ Read-Only Functions

# Returns the L2 public key and game-related player data for a user.
@view
func get_user_info{
        syscall_ptr : felt*,
        pedersen_ptr : HashBuiltin*,
        range_check_ptr
    }(
        user_id : felt,
    ) -> (
        user_data : felt
    ):
    # The GameEngine contract calls this function when a player
    # takes a turn. This ensures a user is allowed to play.
    # The user_data provides different properties during gameplay.
    return user_data.read(user_id)
end

# Returns a 4-bit value at a particular index for item score.
@view
func unpack_score{
        syscall_ptr : felt*,
        pedersen_ptr : HashBuiltin*,
        bitwise_ptr: BitwiseBuiltin*,
        range_check_ptr
    }(
        user_id : felt,
        index : felt
    ) -> (
        score : felt
    ):
    alloc_locals
    # User data is a binary encoded value with alternating
    # 6-bit id followed by a 4-bit score (see top of file).
    let (local data) = user_data.read(user_id)
    local syscall_ptr : felt* = syscall_ptr
    local pedersen_ptr : HashBuiltin* = pedersen_ptr
    local bitwise_ptr: BitwiseBuiltin* = bitwise_ptr
    # 1. Create a 4-bit mask at and to the left of the index
    # E.g., 000111100 = 2**2 + 2**3 + 2**4 + 2**5
    # E.g.,  2**(i) + 2**(i+1) + 2**(i+2) + 2**(i+3) = (2**i)(15)
    let (power) = pow(2, index)
    # 1 + 2 + 4 + 8 = 15
    let mask = 15 * power

    # 2. Apply mask using bitwise operation: mask AND data.
    let (masked) = bitwise_and(mask, data)

    # 3. Shift element right by dividing by the order of the mask.
    let (result, _) = unsigned_div_rem(masked, power)

    # If no score is set for this users item (e.g., registry
    # has not been correctly initialised for this user), give the
    # item a score of 5 (middle-range score).
    local score
    if result == 0:
        assert score = 5
    else:
        assert score = result
    end

    return (score)
end

@view
func get_user_count{
        syscall_ptr : felt*,
        pedersen_ptr : HashBuiltin*,
        range_check_ptr
    }() -> (
        user_count : felt
    ):
    return user_count.read()
end

##### External Functions #####

# User with specific token calls to save their details the game.
@external
func register_user{
        syscall_ptr : felt*,
        pedersen_ptr : HashBuiltin*,
        range_check_ptr
    }(
        data : felt
    ):
    # Performs a check on either:
    # 1) Merkle claim or
    # 2) Ownership of L2-bridged ERC721, ERC1155 or ERC20 token

    # Allocates the user a user_id
    let (user_id) = get_caller_address()

    # Saves the user_id, L2_public_key and user_data

    # Testing
    # ensure user data is non-zero
    assert_not_zero(data)

    # ensure user isn't already registered
    let (existing) = user_data.read(user_id)
    assert existing = 0

    # User data may be a binary encoding of all assets.
    # 00000000000000000000000000010000000010001
    #                            ^ RR         ^ shovel
    user_data.write(user_id, data)

    # Increment user count
    let (count) = user_count.read()
    user_count.write(count + 1)

    return ()
end
