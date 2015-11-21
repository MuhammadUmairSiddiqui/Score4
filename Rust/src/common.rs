pub const WIDTH: u32 = 7;
pub const HEIGHT: u32 = 6;
pub const MAX_DEPTH: i32 = 7;
pub const ORANGE_WINS: i32 = 1000000;
pub const YELLOW_WINS: i32 = -ORANGE_WINS;

pub type Board = [[i32; WIDTH as usize]; HEIGHT as usize];

#[inline(always)]
pub fn other_color(color:i32) -> i32 {
    -color
}

pub fn score_board(board:&Board) -> i32 {
    let mut counts = [0,0,0,0,0,0,0,0,0];

    let myincr = |arr:&mut[i32; 9], idx:i32| {
         arr[idx as usize] = arr[idx as usize] + 1;
    };

    // Horizontal spans
    for yy in 0..HEIGHT {
        let y = yy as usize;
        let mut score: i32 = board[y][0] + board[y][1] + board[y][2];
        for xx in 3..WIDTH {
            let x = xx as usize;
            score = score + board[y][x];
            myincr(&mut counts, score+4);
            score = score - board[y][x-3];
        }
    }

    // Vertical spans
    for xx in 0..WIDTH {
        let x = xx as usize;
        let mut score: i32 = board[0][x] + board[1][x] + board[2][x];
        for yy in 3..HEIGHT {
            let y = yy as usize;
            score = score + board[y][x];
            myincr(&mut counts, score+4);
            score = score - board[y-3][x];
        }
    }

    // Down-right (and up-left) diagonals
    for yy in 0..HEIGHT-3 {
        let y = yy as usize;
        for xx in 0 .. WIDTH-3 {
            let x = xx as usize;
            let mut score: i32 = 0;
            for idx in 0 .. 4 {
                score = score + board[y+idx][x+idx];
            }
            myincr(&mut counts, score+4);
        }
    }

    // up-right (and down-left) diagonals
    for yy in 3..HEIGHT {
        let y = yy as usize;
        for xx in 0..WIDTH-3 {
            let x = xx as usize;
            let mut score: i32 = 0;
            for idx in 0..4 {
                score = score + board[y-idx][x+idx];
            }
            myincr(&mut counts, score+4);
        }
    }

    if counts[0] != 0 {
        YELLOW_WINS
    } else if counts[8] != 0 {
        ORANGE_WINS
    } else {
        counts[5] + 2*counts[6] + 5*counts[7] -
            counts[3] - 2*counts[2] - 5*counts[1]
    }
}

/* vim: set expandtab ts=8 sts=4 shiftwidth=4 */
