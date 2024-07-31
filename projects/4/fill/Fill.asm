// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed,
// the screen should be cleared.

// SCREEN IS 256 ROWS BY 512
// @SCREEN REFERS TO 16384
// 8K MEMORY BLOCK OF 16-BIT WORDS STARTING AT 16384
// 32x16=512, 32 16 bit words per row

// 0 IS SET, 1 IS NOT SET
    @WHITEFLAG //0
    M=0 // SCREEN IS WHITE BY DEFAULT
    @BLACKFLAG
    M=1
    @255
    D=A
    @MAXROW
    M=D
    @31
    D=A
    @MAXCOL
    M=D //11

(LOOP)
    // if keyboard input !=0 (INPUT) jump to black (blacken screen)
    @KBD //12
    D=M
    @BLACK
    D;JGT

    // check white flag, if >0 go to WHITE otherwise loop
    @WHITEFLAG
    D=M
    @WHITE
    D;JGT

    @LOOP
    0;JMP


// IF BLACKFLAG 0, JUMP BACK TO LOOK, OTHERWISE DRAW THEN SET TO 0 AND SET WHI
// TEFLAG TO 1
(BLACK)
    @BLACKFLAG //22
    D=M
    @LOOP
    D;JEQ // IF BLACKFLAG IS SET, LOOP AGAIN

    // INITALIZE VARIABLES
    @ROW //26
    M=0
    @ROWMEM
    M=0
    @BLACKFLAG
    M=0 // SET
    @WHITEFLAG
    M=1 //UNSET

    @BLOOP //34
    0;JMP //SKIP RESET COL COUNTER LOGIC/ROW

    (RESETCOL_ADDROW) //RESET AFTER EACH ROW
    @COL
    M=0
    @ROW
    M=M+1
    @MAXCOL
    D=M+1
    @ROWMEM
    M=D+M

    (BLOOP) //MAIN LOGIC TO SET SCREEN TO BLACK
    @ROW //44
    D=M
    @MAXROW
    D=D-M
    @LOOP   // EXIT BLACKLOOP BACK TO MAIN LOOP
    D;JGT //ROW > MAXROW

    @COL //50
    D=M
    @MAXCOL
    D=D-M
    @RESETCOL_ADDROW //GO ONTO NEXT COLUMN
    D;JGT //COL > MAXCOL

    // SCREEN 16BITS TO BLACK
    // SCREEN + ROW*32 + COL/16
    // ROW*32 -  ROWCOUNTER - ROWCOUNTER = 32 + 32...ROW TIMES
    // COL/16 - DOESN'T MATTER WE ONLY CARE ABOUT REACHING EACH 16BIT REGISTER
    // IN EACH ROW, SO 0...31 COUNTER
    @ROWMEM //56
    D=M
    @COL
    D=D+M
    @SCREEN
    D=D+A
    @SCREEN16BIT
    A=D
    M=-1
    @COL
    M=M+1
    @BLOOP
    0;JMP //66

(WHITE)
    // INITALIZE VARIABLES
    @ROW
    M=0
    @ROWMEM
    M=0
    @WHITEFLAG
    M=0 //SET
    @BLACKFLAG
    M=1 //UNSET

    @WLOOP
    0;JMP //SKIP RESET COL COUNTER LOGIC/ROW

    (RESETCOL_ADDROWW) //RESET AFTER EACH ROW
    @COL
    M=0
    @ROW
    M=M+1
    @MAXCOL
    D=M+1
    @ROWMEM
    M=D+M

    (WLOOP) //MAIN LOGIC TO SET SCREEN TO WHITE
    @ROW
    D=M
    @MAXROW
    D=D-M
    @LOOP   // EXIT WHITELOOP BACK TO MAIN LOOP
    D;JGT //ROW > MAXROW

    @COL
    D=M
    @MAXCOL
    D=D-M
    @RESETCOL_ADDROWW //GO ONTO NEXT COLUMN
    D;JGT //COL > MAXCOL

    @ROWMEM
    D=M
    @COL
    D=D+M
    @SCREEN
    D=D+A
    @SCREEN16BIT
    A=D
    M=0
    @COL
    M=M+1
    @WLOOP
    0;JMP

//loop until whole screen is black, then copy code but with 0's for white
//then set white flag after whiten screen and the program is done



//PSEUDO CODE FOR SETTING SCREEN TO BLACK
//ROW 0-511
//COL 0-31
//RAM[SCREEN + ROW*32 +COL/16]
//COL%16=POSITION (DOESN'T MATTER FOR THIS PROGRAM)
//
//ROW=0
//MAXROW=511
//MAXCOL=31
//GO TO LOOP // SKIP RESET FOR FIRST RUN
//
//(RESETCOL_ADDROW)
//COL=0
//ROW+=1
//
//(LOOP)
//IF ROW > MAXROW GO TO END //ROW - MAXROW > 0
//IF COL > MAXCOL GO TO RESETCOL_ADDROW //COL - MAXCOL > 0
//16BIT_PIXEL_LOC = SCREEN + ROW*32 + COL/16
//16BIT_PIXEL_LOC = 1
//GO TO LOOP

