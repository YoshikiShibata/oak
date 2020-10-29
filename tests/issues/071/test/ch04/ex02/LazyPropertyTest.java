/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
package ch04.ex02;

import javafx.beans.property.Property;
import javafx.beans.property.SimpleStringProperty;
import org.junit.Test;
import static org.junit.Assert.*;

/**
 *
 * Test code for LazyProperty
 */
public class LazyPropertyTest {
    
    @Test(expected=NullPointerException.class)
    public void testNullPointerException() {
        new LazyProperty<Object>(null);
    }
    
    @Test
    public void simpleSetAndGet() {
        LazyProperty<String> lp = new LazyProperty<>(SimpleStringProperty::new);
        
        lp.setValue("hello");
        assertEquals("hello", lp.getValue());
    }
    
    @Test
    public void testProperty() {
        LazyProperty<String> lp = new LazyProperty<>(SimpleStringProperty::new);
        
        lp.setValue("hello");
        Property<String> p = lp.property();
        assertEquals("hello", lp.getValue());
    }
}
