/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
package ch04.ex02;

import java.util.Objects;
import java.util.function.Supplier;
import javafx.beans.property.Property;

/**
 * Consider a class with many JavaFX properties, such as a chart or table.
 * Chances are that in a particular application, most properties never have
 * listeners attached to them. It is therefore wasteful to have a property
 * object per property. Show how the property can be set up on demand, first
 * using a regular field for storing the property value, and then using a
 * property object only when the xxxProperty() method is called for the first
 * time.
 *
 * チャートやテーブルといった多くのJavaFX プロパティを持つクラスを考えなさい。特定
 * のアプリケーションでは、ほとんどのプロパティには決してリスナーが登録されない可
 * 能性が高いです。したがって、プロパティごとにプロパティオブジェクトを持つことは
 * 無駄です。プロパティ値を保存するために最初に普通のフィールドを使用して、初めて xxxProperty()
 * メソッドが呼び出されたときにだけプロパティオブジェクトを使用す るように、要求に応じてプロパティを構築する方法を示しなさい。
 */
public class LazyProperty<T> {

    private T value = null;
    private Property<T> property = null;
    private final Supplier<Property<T>> supplier;

    /**
     * Constructs with a Supplier which is used to create an instance of
     * Property<T>.
     * 
     * @param supplier supplier to create an instance.
     * @throws NullPointerException if supplier is null.
     */
    public LazyProperty(Supplier<Property<T>> supplier) {
        Objects.requireNonNull(supplier, "supplier is null");
        this.supplier = supplier;
    }

    public final void setValue(T value) {
        if (property != null) {
            property.setValue(value);
        } else {
            this.value = value;
        }
    }

    public final T getValue() {
        return property != null ? property.getValue() : value;
    }

    public final Property<T> property() {

        if (property != null) {
            return property;
        }

        property = supplier.get();
        property.setValue(value);
        value = null;
        return property;
    }

}
