
'
' Օրինակ 8. փոփոխականների տիրույթներ
'
SUB Main
    LET a = 1

    LET c = a + 2
    LET d = c + e  ' սխալ, e-ն արժեք չունի

    FOR i = 1 TO 10
        PRINT i
        IF i \ 2 = 0 THEN
            LET c = i * 2 + 1
            PRINT c
        END IF
        PRINT c  ' սա 8-րդ տողի c-ն է
    END FOR

    PRINT i  ' սխալ, i-ն սահմանված չէ

END SUB
