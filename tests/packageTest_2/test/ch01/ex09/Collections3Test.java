/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

package ch01.ex09;

import java.util.ArrayList;
import org.junit.Test;

/**
 *
 * @author yoshiki
 */
public class Collections3Test {
    
    static class ArrayList2<E> extends ArrayList<E> implements Collection2<E> {
        
    }
    
    @Test 
    public void forEachIf() {
        ArrayList2<Integer> list = new ArrayList2<>();
        for (int i = 0; i < 100; i++)
            list.add(i);
        
        list.forEachIf(
                i -> {System.out.println(i);}, 
                i -> { return i == 10;} );
        
        list.forEach(i -> {
           if (i == 10)
               System.out.println(i);
        });
    }
}
